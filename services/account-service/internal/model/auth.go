package model

import (
	"time"

	"github.com/google/uuid"
)

type SessionData struct {
	UserID                uuid.UUID     `json:"uid"`
	DeviceID              string        `json:"did"`
	DataSession           string        `json:"-"`
	RefreshToken          string        `json:"-"`
	SessionExpiration     time.Duration `json:"-"`
	HashedToken           string        `json:"-"`
	HashedTokenExpiration time.Duration `json:"-"`
}

type SessionPayload struct {
	UserID   uuid.UUID `json:"uid"`
	DeviceID string    `json:"did"`
}

type GoogleLoginRequest struct {
	Token    string `json:"token"`
	DeviceID string `json:"-"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"` // get from cookie
	DeviceID     string `json:"-"`
}

type LogoutRequest struct {
	RefreshToken string    `json:"refresh_token"` // get from cookie
	UserID       uuid.UUID `json:"-"`
	DeviceID     string    `json:"-"`
}

type RevokeAllSessionsRequest struct {
	AccessToken string    `json:"-"`
	UserID      uuid.UUID `json:"-"`
	DeviceID    string    `json:"-"`
}

type AuthUserResponse struct {
	ID         string  `json:"id"`
	Email      string  `json:"email"`
	FirstName  *string `json:"first_name,omitempty"`
	LastName   *string `json:"last_name,omitempty"`
	AvatarURL  *string `json:"avatar_url,omitempty"`
	IsVerified bool    `json:"is_verified"`
}

type AuthTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	// User         AuthUserResponse `json:"user"`
}
