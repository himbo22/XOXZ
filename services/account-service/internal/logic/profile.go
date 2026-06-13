package logic

import (
	"context"
	"time"

	"github.com/google/uuid"
	mediaGrpc "github.com/himbo22/xoxz/account-service/internal/adapter/grpc"
	"github.com/himbo22/xoxz/account-service/internal/config"
	_const "github.com/himbo22/xoxz/account-service/internal/const"
	"github.com/himbo22/xoxz/account-service/internal/domain/repository"
	"github.com/himbo22/xoxz/account-service/internal/model"
	"github.com/himbo22/xoxz/account-service/internal/util"
	"github.com/himbo22/xoxz/common-service/protobuf/media"
	xoxz "github.com/himbo22/xoxz/common-service/xoxz/logger"
)

type ProfileLogic struct {
	config      *config.Config
	mediaClient mediaGrpc.MediaClient
	logger      xoxz.XoxzLogger
	userRepo    repository.UserRepository
}

func NewProfileLogic(
	config *config.Config,
	logger xoxz.XoxzLogger,
	mediaClient mediaGrpc.MediaClient,
	userRepo repository.UserRepository,
) *ProfileLogic {
	return &ProfileLogic{
		config:      config,
		logger:      logger,
		mediaClient: mediaClient,
		userRepo:    userRepo,
	}
}

func (p *ProfileLogic) GetProfile(ctx context.Context, userID uuid.UUID) (model.ProfileResponse, error) {
	user, err := p.userRepo.FindByID(ctx, userID)
	if err != nil {
		return model.ProfileResponse{}, util.NewErrorByCode(_const.CodeInternalError, "database error getting user profile")
	}
	if user == nil {
		return model.ProfileResponse{}, util.NewErrorByCode(_const.CodeUserNotFound)
	}

	return model.ProfileResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
		AvatarURL: user.AvatarURL,
		Bio:       user.Bio,
		Status:    user.Status,
	}, nil
}

func (p *ProfileLogic) GetPublicProfile(ctx context.Context, username string) (model.PublicProfileResponse, error) {
	user, err := p.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return model.PublicProfileResponse{}, util.NewErrorByCode(_const.CodeInternalError, "database error getting public profile")
	}
	if user == nil {
		return model.PublicProfileResponse{}, util.NewErrorByCode(_const.CodeUserNotFound)
	}

	return model.PublicProfileResponse{
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		AvatarURL: user.AvatarURL,
		Status:    user.Status,
		Bio:       user.Bio,
	}, nil
}

func (p *ProfileLogic) UpdateProfile(ctx context.Context, userID uuid.UUID, req model.UpdateProfileRequest) (model.ProfileResponse, error) {
	user, err := p.userRepo.FindByID(ctx, userID)
	if err != nil {
		return model.ProfileResponse{}, util.NewErrorByCode(_const.CodeInternalError, "database error getting user profile")
	}
	if user == nil {
		return model.ProfileResponse{}, util.NewErrorByCode(_const.CodeUserNotFound)
	}

	if req.FirstName != nil {
		user.FirstName = req.FirstName
	}
	if req.Username != nil {
		existedUser, err := p.userRepo.FindByUsername(ctx, *req.Username)
		if err != nil {
			return model.ProfileResponse{}, util.NewErrorByCode(_const.CodeInternalError, "database error checking username")
		}
		if existedUser != nil && existedUser.ID != user.ID {
			return model.ProfileResponse{}, util.NewErrorByCode(_const.CodeInvalidRequest, "username already exists")
		}

		user.Username = req.Username
	}
	if req.LastName != nil {
		user.LastName = req.LastName
	}
	if req.Phone != nil {
		user.Phone = req.Phone
	}
	if req.Bio != nil {
		user.Bio = req.Bio
	}

	if err := p.userRepo.Update(ctx, user); err != nil {
		return model.ProfileResponse{}, util.NewErrorByCode(_const.CodeInternalError, "database error updating user profile")
	}

	return model.ProfileResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
		AvatarURL: user.AvatarURL,
		Bio:       user.Bio,
		Status:    user.Status,
	}, nil
}

func (p *ProfileLogic) UpdateAvatar(ctx context.Context, req model.UpdateAvatarRequest) (model.UpdateAvatarResponse, error) {
	// check user first - internal data first
	user, err := p.userRepo.FindByID(ctx, req.UserID)
	if err != nil {
		return model.UpdateAvatarResponse{}, err
	}
	if user == nil {
		// not found
		return model.UpdateAvatarResponse{}, util.NewErrorByCode(_const.CodeUserNotFound)
	}

	payload := mediaGrpc.ToCommitFileRequest(req)

	response, err := p.mediaClient.CommitFile(ctx, payload)
	if err != nil {
		return model.UpdateAvatarResponse{}, err
	}
	if !response.GetSuccess() {
		return model.UpdateAvatarResponse{}, util.NewErrorByCode(_const.CodeStorageError)
	}
	// update avatar

	user.AvatarURL = &req.PerPath
	if err := p.userRepo.Update(ctx, user); err != nil {
		// delete file in media-svc
		p.mediaClient.DeleteFile(ctx, &media.DeleteFileRequest{ObjectPath: req.PerPath})
		// Give rollback an independent 5-second timeout.
		rollbackCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Delete the file through gRPC.
		_, delErr := p.mediaClient.DeleteFile(rollbackCtx, &media.DeleteFileRequest{ObjectPath: req.PerPath})

		if delErr != nil {
			// File deletion is not guaranteed, so log the orphaned object.
			p.logger.Errorf("CRITICAL: Orphaned object left in MinIO!", xoxz.String("path", req.PerPath), xoxz.Error(delErr))
		}
		return model.UpdateAvatarResponse{}, err
	}

	return model.UpdateAvatarResponse{NewAvatar: req.PerPath}, nil
}
