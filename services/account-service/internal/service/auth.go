package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
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

type AuthService interface {
	AuthenticateWithGoogle(ctx context.Context, req model.GoogleLoginRequest) (res model.AuthTokenResponse, err error)
	RefreshToken(ctx context.Context, req model.RefreshTokenRequest) (model.AuthTokenResponse, error)
	Logout(ctx context.Context, req model.LogoutRequest) error
	RevokeAllSessions(ctx context.Context, req model.RevokeAllSessionsRequest) error
}

type GoogleOAuthProvider interface {
	VerifyGoogleToken(ctx context.Context, token, clientID string) (*model.TokenPayload, error)
}

type authService struct {
	config             *config.Config
	logger             xoxz.XoxzLogger
	userRepo           repository.UserRepository
	identityRepo       repository.IdentityRepository
	roleRepository     repository.RoleRepository
	userRoleRepository repository.UserRoleRepository
	txRunner           repository.TxRunner
	googleAuthProvider GoogleOAuthProvider
	redisRepo          repository.RedisRepository
}

func NewAuthService(
	config *config.Config,
	logger xoxz.XoxzLogger,
	userRepo repository.UserRepository,
	identityRepo repository.IdentityRepository,
	roleRepository repository.RoleRepository,
	userRoleRepository repository.UserRoleRepository,
	txRunner repository.TxRunner,
	googleAuthProvider GoogleOAuthProvider,
	redisRepo repository.RedisRepository,
) AuthService {
	return &authService{
		config:             config,
		logger:             logger,
		userRepo:           userRepo,
		identityRepo:       identityRepo,
		roleRepository:     roleRepository,
		userRoleRepository: userRoleRepository,
		txRunner:           txRunner,
		googleAuthProvider: googleAuthProvider,
		redisRepo:          redisRepo,
	}
}

func (a *authService) AuthenticateWithGoogle(ctx context.Context, req model.GoogleLoginRequest) (res model.AuthTokenResponse, err error) {
	ctx, end := telemetry.StartSpan(ctx, "test1", "service")
	defer func() { end(err) }()

	var payload *model.TokenPayload
	payload, err = a.verifyAndValidateToken(ctx, req.Token)
	if err != nil {
		return
	}

	hashedToken := util.HashToken(req.Token)
	if err = a.preventReplayAttack(ctx, hashedToken); err != nil {
		return
	}

	var identity *entity.Identity
	identity, err = a.identityRepo.FindBySubjectAndProvider(ctx, payload.Subject, _const.ProviderGoogle)
	if err != nil {
		return
	}
	if identity == nil {
		res, err = a.registerWithGoogle(ctx, payload, req.DeviceID, hashedToken)
		return
	}

	err = a.checkUserStatus(ctx, identity.UserID)
	if err != nil {
		return
	}

	return a.loginWithGoogle(ctx, payload, identity, req.DeviceID, hashedToken)
}

func (a *authService) verifyAndValidateToken(ctx context.Context, token string) (*model.TokenPayload, error) {
	payload, err := a.googleAuthProvider.VerifyGoogleToken(ctx, token, a.config.Auth.GoogleClientId)
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

func (a *authService) preventReplayAttack(ctx context.Context, token string) error {
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

func (a *authService) checkUserStatus(ctx context.Context, userID uuid.UUID) error {
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

func (a *authService) RefreshToken(ctx context.Context, req model.RefreshTokenRequest) (model.AuthTokenResponse, error) {
	sessionKey := _const.RefreshTokenKey(req.RefreshToken)
	sessionJSON, err := a.redisRepo.Get(ctx, sessionKey)
	if err != nil {
		return model.AuthTokenResponse{}, util.NewErrorByCode(_const.CodeInternalError)
	}
	if sessionJSON == "" {
		return model.AuthTokenResponse{}, util.NewErrorByCode(_const.CodeExpiredRefreshToken)
	}

	var payload model.SessionPayload
	if err := json.Unmarshal([]byte(sessionJSON), &payload); err != nil {
		return model.AuthTokenResponse{}, err
	}

	if payload.DeviceID != req.DeviceID {
		return model.AuthTokenResponse{}, util.NewErrorByCode(_const.CodeUnauthorized)
	}

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

func (a *authService) Logout(ctx context.Context, req model.LogoutRequest) error {
	if err := a.redisRepo.RemoveSession(ctx, req.UserID, req.DeviceID, req.RefreshToken); err != nil {
		return err
	}

	return nil
}

func (a *authService) RevokeAllSessions(ctx context.Context, req model.RevokeAllSessionsRequest) error {
	if err := a.redisRepo.RevokeAllSessions(ctx, req.UserID); err != nil {
		return err
	}

	return nil
}

func (a *authService) authTokenResponse(ctx context.Context, userID uuid.UUID, deviceID string, hashedToken string, hashedTokenTTL int64) (model.AuthTokenResponse, error) {
	at, err := util.GenerateAccessToken(config.AppPrivateKey, userID, deviceID, config.AccessTokenExpiryTime)
	if err != nil {
		return model.AuthTokenResponse{}, err
	}

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

func (a *authService) loginWithGoogle(ctx context.Context, payload *model.TokenPayload, identity *entity.Identity, deviceID string, hashedToken string) (res model.AuthTokenResponse, err error) {
	err = a.txRunner.RunInTx(ctx, func(txCtx context.Context) error {
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

func (a *authService) registerWithGoogle(ctx context.Context, payload *model.TokenPayload, deviceID string, hashedToken string) (res model.AuthTokenResponse, err error) {
	var userID uuid.UUID
	userID, err = uuid.NewV7()
	if err != nil {
		return
	}

	err = a.txRunner.RunInTx(ctx, func(txCtx context.Context) error {
		role, err := a.roleRepository.FindByName(ctx, _const.ROLE_USER)
		if err != nil {
			return util.NewErrorByCode(_const.CodeInternalError, "database error getting role")
		}
		if role == nil {
			return util.NewErrorByCode(_const.CodeInternalError, "database error getting role")
		}

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

		userRole := &entity.UserRole{
			RoleID: role.ID,
			UserID: userID,
		}

		if err = a.userRoleRepository.Create(txCtx, userRole); err != nil {
			return err
		}

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
