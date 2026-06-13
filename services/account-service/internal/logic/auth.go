package logic

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	adapter_oauth "github.com/himbo22/xoxz/account-service/internal/adapter/oauth"
	"github.com/himbo22/xoxz/account-service/internal/config"
	_const "github.com/himbo22/xoxz/account-service/internal/const"
	"github.com/himbo22/xoxz/account-service/internal/domain/entity"
	"github.com/himbo22/xoxz/account-service/internal/domain/repository"
	"github.com/himbo22/xoxz/account-service/internal/model"
	"github.com/himbo22/xoxz/account-service/internal/util"
	"github.com/himbo22/xoxz/common-service/monitoring/telemetry"
	xoxz "github.com/himbo22/xoxz/common-service/xoxz/logger"
	"gorm.io/datatypes"
)

type AuthLogic struct {
	config             *config.Config
	logger             xoxz.XoxzLogger
	userRepo           repository.UserRepository
	identityRepo       repository.IdentityRepository
	roleRepository     repository.RoleRepository
	userRoleRepository repository.UserRoleRepository
	txRunner           repository.TxRunner
	redisRepo          repository.RedisRepository
}

func NewAuthLogic(
	config *config.Config,
	logger xoxz.XoxzLogger,
	userRepo repository.UserRepository,
	identityRepo repository.IdentityRepository,
	roleRepository repository.RoleRepository,
	userRoleRepository repository.UserRoleRepository,
	txRunner repository.TxRunner,
	redisRepo repository.RedisRepository,
) *AuthLogic {
	return &AuthLogic{
		config:             config,
		logger:             logger,
		userRepo:           userRepo,
		identityRepo:       identityRepo,
		roleRepository:     roleRepository,
		userRoleRepository: userRoleRepository,
		txRunner:           txRunner,
		redisRepo:          redisRepo,
	}
}

func (a *AuthLogic) AuthenticateWithGoogle(ctx context.Context, req model.GoogleLoginRequest) (res model.AuthTokenResponse, err error) {
	ctx, end := telemetry.StartSpan(ctx, "test1", "logic-1")
	defer func() { end(err) }()

	// 1. verify payload
	var payload *model.TokenPayload
	payload, err = a.verifyAndValidateToken(ctx, req.Token)
	if err != nil {
		return
	}

	// 2. check if token exists in redis -> prevent reply attack
	// we hash token here because google token is very long so we need to do it shorter
	hashedToken := util.HashToken(req.Token)
	if err = a.preventReplayAttack(ctx, hashedToken); err != nil {
		return
	}

	// 3. Check if user is in system or not
	var identity *entity.Identity
	identity, err = a.identityRepo.FindBySubjectAndProvider(ctx, payload.Subject, _const.ProviderGoogle)
	if err != nil {
		return
	}
	// 3.1 register
	if identity == nil {
		res, err = a.registerWithGoogle(ctx, payload, req.DeviceID, hashedToken)
		return
	}

	// 4. check user status : active/lock
	err = a.checkUserStatus(ctx, identity.UserID)
	if err != nil {
		return
	}

	// login
	return a.logicWithGoogle(ctx, payload, identity, req.DeviceID, hashedToken)
}

func (a *AuthLogic) verifyAndValidateToken(ctx context.Context, token string) (*model.TokenPayload, error) {
	payload, err := adapter_oauth.VerifyGoogleToken(ctx, token, a.config.Auth.GoogleClientId)
	if err != nil {
		a.logger.Error("google token verification failed", xoxz.Error(err))
		return nil, util.NewErrorByCode(_const.CodeInvalidGoogleToken)
	}

	if !util.GetBool(payload.Claims, "email_verified", false) {
		return nil, util.NewErrorByCode(_const.CodeInvalidGoogleToken, "email is not verified")
	}

	if payload.Audience != a.config.Auth.GoogleClientId {
		return nil, util.NewErrorByCode(_const.CodeInvalidGoogleToken, "issuer audience")
	}

	if time.Now().Unix() > payload.Expires {
		return nil, util.NewErrorByCode(_const.CodeInvalidGoogleToken, "expires token")
	}

	if payload.Issuer != "https://accounts.google.com" {
		return nil, util.NewErrorByCode(_const.CodeInvalidGoogleToken, "invalid issuer")
	}

	return payload, nil
}

func (a *AuthLogic) preventReplayAttack(ctx context.Context, token string) error {
	googleTokenKey := _const.GoogleTokenKey(token)

	isExisted, err := a.redisRepo.Exists(ctx, googleTokenKey)
	if err != nil {
		return err
	}
	if isExisted {
		return util.NewErrorByCode(_const.CodeUnauthorized, "token already used")
	}

	return nil
}

func (a *AuthLogic) checkUserStatus(ctx context.Context, userID uuid.UUID) error {
	user, err := a.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return util.NewErrorByCode(_const.CodeUserNotFound)
	}
	if user.Status != nil && *user.Status != "ACTIVE" {
		return util.NewErrorByCode(_const.CodeForbidden, "user is locked")
	}
	return nil
}

func (a *AuthLogic) RefreshToken(ctx context.Context, req model.RefreshTokenRequest) (model.AuthTokenResponse, error) {
	// check authenticated token
	sessionKey := _const.RefreshTokenKey(req.RefreshToken)
	sessionJSON, err := a.redisRepo.Get(ctx, sessionKey)
	if err != nil {
		return model.AuthTokenResponse{}, util.NewErrorByCode(_const.CodeInternalError)
	}
	if sessionJSON == "" {
		return model.AuthTokenResponse{}, util.NewErrorByCode(_const.CodeExpiredRefreshToken)
	}

	// parse json -> struct
	var payload model.SessionPayload
	if err := json.Unmarshal([]byte(sessionJSON), &payload); err != nil {
		return model.AuthTokenResponse{}, err
	}

	if payload.DeviceID != req.DeviceID {
		return model.AuthTokenResponse{}, util.NewErrorByCode(_const.CodeUnauthorized)
	}

	// check user status & device status
	user, err := a.userRepo.FindByID(ctx, payload.UserID)
	if err != nil {
		return model.AuthTokenResponse{}, util.NewErrorByCode(_const.CodeInternalError)
	}
	if user == nil {
		return model.AuthTokenResponse{}, util.NewErrorByCode(_const.CodeUserNotFound)
	}
	if user.Status != nil && *user.Status != "ACTIVE" {
		return model.AuthTokenResponse{}, util.NewErrorByCode(_const.CodeForbidden, "user is locked")
	}

	// generate new AT & RT
	newAccessToken, err := util.GenerateAccessToken(config.AppPrivateKey, payload.UserID, payload.DeviceID, config.AccessTokenExpiryTime)
	if err != nil {
		return model.AuthTokenResponse{}, err
	}
	newRefreshToken := util.GenerateRefreshToken()

	rtValueBytes, _ := json.Marshal(payload)

	sessionData := model.SessionData{
		UserID:            payload.UserID,
		DeviceID:          payload.DeviceID,
		RefreshToken:      newRefreshToken,
		DataSession:       string(rtValueBytes),
		SessionExpiration: config.RefreshTokenExpiryTime,
	}
	if err := a.redisRepo.CreateSession(ctx, sessionData); err != nil {
		return model.AuthTokenResponse{}, err
	}

	return model.AuthTokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (a *AuthLogic) Logout(ctx context.Context, req model.LogoutRequest) error {
	if err := a.redisRepo.RemoveSession(ctx, req.UserID, req.DeviceID, req.RefreshToken); err != nil {
		return err
	}

	return nil
}

func (a *AuthLogic) RevokeAllSessions(ctx context.Context, req model.RevokeAllSessionsRequest) error {
	if err := a.redisRepo.RevokeAllSessions(ctx, req.UserID); err != nil {
		return err
	}

	return nil
}

func (a *AuthLogic) authTokenResponse(ctx context.Context, userID uuid.UUID, deviceID string, hashedToken string, hashedTokenTTL int64) (model.AuthTokenResponse, error) {
	at, err := util.GenerateAccessToken(config.AppPrivateKey, userID, deviceID, config.AccessTokenExpiryTime)
	if err != nil {
		return model.AuthTokenResponse{}, err
	}

	ctx, end := telemetry.StartSpan(ctx, "test1", "create token")
	defer func() { end(err) }()

	refreshToken := util.GenerateRefreshToken()

	sessionPayload := model.SessionPayload{
		UserID:   userID,
		DeviceID: deviceID,
	}

	rtValue, err := json.Marshal(sessionPayload)
	if err != nil {
		return model.AuthTokenResponse{}, err
	}
	jsonStr := string(rtValue)

	sessionData := model.SessionData{
		UserID:                userID,
		DeviceID:              deviceID,
		RefreshToken:          refreshToken,
		DataSession:           jsonStr,
		HashedToken:           hashedToken,
		SessionExpiration:     config.RefreshTokenExpiryTime,
		HashedTokenExpiration: util.UnixToDuration(hashedTokenTTL),
	}

	if err := a.redisRepo.CreateGoogleSession(
		ctx,
		sessionData,
	); err != nil {
		return model.AuthTokenResponse{}, err
	}

	res := model.AuthTokenResponse{
		AccessToken:  at,
		RefreshToken: refreshToken,
	}

	return res, nil
}

func (a *AuthLogic) logicWithGoogle(ctx context.Context, payload *model.TokenPayload, identity *entity.Identity, deviceID string, hashedToken string) (res model.AuthTokenResponse, err error) {
	ctx, end := telemetry.StartSpan(ctx, "test1", "logic-login")
	defer func() { end(err) }()

	err = a.txRunner.RunInTx(ctx, func(txCtx context.Context) error {
		// update user last logic
		if err := a.userRepo.UpdateLastLogin(txCtx, identity.UserID); err != nil {
			return util.NewErrorByCode(_const.CodeInternalError, "database error updating user")
		}
		rawClaimsBytes, err := json.Marshal(payload.Claims)
		if err != nil {
			return util.NewErrorByCode(_const.CodeInternalError, "internal data error")
		}

		identity.ProviderData = rawClaimsBytes

		if err := a.identityRepo.Upsert(txCtx, identity); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return
	}

	return a.authTokenResponse(ctx, identity.UserID, deviceID, hashedToken, payload.Expires)
}

func (a *AuthLogic) registerWithGoogle(ctx context.Context, payload *model.TokenPayload, deviceID string, hashedToken string) (res model.AuthTokenResponse, err error) {
	var userID uuid.UUID
	userID, err = uuid.NewV7()
	if err != nil {
		return
	}

	// create user + create identity
	ctx, end := telemetry.StartSpan(ctx, "test1", "register")
	defer func() { end(err) }()

	err = a.txRunner.RunInTx(ctx, func(txCtx context.Context) error {
		// TODO: role
		role, err := a.roleRepository.FindByName(ctx, _const.ROLE_USER)
		if err != nil {
			return util.NewErrorByCode(_const.CodeInternalError, "database error getting role")
		}
		if role == nil {
			return util.NewErrorByCode(_const.CodeInternalError, "database error getting role")
		}

		// user
		user := &entity.User{
			ID:        userID,
			FirstName: util.GetString(payload.Claims, "given_name", ""),
			LastName:  util.GetString(payload.Claims, "family_name", ""),
			AvatarURL: util.GetString(payload.Claims, "picture", ""),
			Status:    new(_const.USER_STATUS_ACTIVE),
		}
		if err := a.userRepo.Create(txCtx, user); err != nil {
			return err
		}

		// add user_role
		userRole := &entity.UserRole{
			RoleID: role.ID,
			UserID: userID,
		}

		if err = a.userRoleRepository.Create(txCtx, userRole); err != nil {
			return err
		}

		// identity
		identityID, err := uuid.NewV7()
		if err != nil {
			return err
		}
		rawClaimsBytes, err := json.Marshal(payload.Claims)
		if err != nil {
			return err
		}
		identity := &entity.Identity{
			ID:             identityID,
			UserID:         userID,
			Provider:       _const.ProviderGoogle,
			ProviderUserID: payload.Subject,
			ProviderData:   datatypes.JSON(rawClaimsBytes),
		}

		if err := a.identityRepo.Upsert(txCtx, identity); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return
	}

	return a.authTokenResponse(ctx, userID, deviceID, hashedToken, payload.Expires)
}
