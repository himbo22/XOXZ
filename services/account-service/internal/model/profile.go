package model

import (
	"github.com/google/uuid"
)

type ProfileResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  *string   `json:"username,omitempty"`
	Email     *string   `json:"email"`
	FirstName *string   `json:"first_name,omitempty"`
	LastName  *string   `json:"last_name,omitempty"`
	Phone     *string   `json:"phone,omitempty"`
	AvatarURL *string   `json:"avatar_url,omitempty"`
	Bio       *string   `json:"bio,omitempty"`
	Status    *string   `json:"status,omitempty"`
}

type PublicProfileResponse struct {
	Username  *string `json:"username,omitempty"`
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
	Bio       *string `json:"bio,omitempty"`
	Status    *string `json:"status,omitempty"`
}

type UpdateProfileRequest struct {
	Username  *string `json:"username,omitempty"`
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Phone     *string `json:"phone,omitempty"`
	Bio       *string `json:"bio,omitempty"`
}

type UpdateAvatarRequest struct {
	TmpPath string    `json:"tmp_path,omitempty"`
	PerPath string    `json:"per_path,omitempty"`
	UserID  uuid.UUID `json:"-"`
}

type UpdateAvatarResponse struct {
	NewAvatar string `json:"new_avt"`
}

type ProfileResponseDoc struct {
	StatusCode int             `json:"status_code"`
	Message    string          `json:"message"`
	Error      *string         `json:"error,omitempty"`
	Data       ProfileResponse `json:"data"`
}

type PublicProfileResponseDoc struct {
	StatusCode int                   `json:"status_code"`
	Message    string                `json:"message"`
	Error      *string               `json:"error,omitempty"`
	Data       PublicProfileResponse `json:"data"`
}

type UpdateAvatarResponseDoc struct {
	StatusCode int                  `json:"status_code"`
	Message    string               `json:"message"`
	Error      *string              `json:"error,omitempty"`
	Data       UpdateAvatarResponse `json:"data"`
}
