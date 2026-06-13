package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/himbo22/xoxz/account-service/internal/logic"
	"github.com/himbo22/xoxz/account-service/internal/model"
)

type ProfileService interface {
	GetProfile(ctx context.Context, userID uuid.UUID) (model.ProfileResponse, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, req model.UpdateProfileRequest) (model.ProfileResponse, error)
	UpdateAvatar(ctx context.Context, req model.UpdateAvatarRequest) (model.UpdateAvatarResponse, error)

	// public profile
	GetPublicProfile(ctx context.Context, username string) (model.PublicProfileResponse, error)
}

type profileService struct {
	profileLogic *logic.ProfileLogic
}

func NewProfileService(profileLogic *logic.ProfileLogic) ProfileService {
	return &profileService{
		profileLogic: profileLogic,
	}
}

// GetProfile implements [ProfileService].
func (p *profileService) GetProfile(ctx context.Context, userID uuid.UUID) (model.ProfileResponse, error) {
	return p.profileLogic.GetProfile(ctx, userID)
}

// GetPublicProfile implements [ProfileService].
func (p *profileService) GetPublicProfile(ctx context.Context, username string) (model.PublicProfileResponse, error) {
	return p.profileLogic.GetPublicProfile(ctx, username)
}

// UpdateAvatar implements [ProfileService].
func (p *profileService) UpdateAvatar(ctx context.Context, req model.UpdateAvatarRequest) (model.UpdateAvatarResponse, error) {
	return p.profileLogic.UpdateAvatar(ctx, req)
}

// UpdateProfile implements [ProfileService].
func (p *profileService) UpdateProfile(ctx context.Context, userID uuid.UUID, req model.UpdateProfileRequest) (model.ProfileResponse, error) {
	return p.profileLogic.UpdateProfile(ctx, userID, req)
}
