package model

import "time"

type CreateArtistAccountRequest struct {
	Email     string `json:"email"`
	StageName string `json:"stage_name"`
}

type CreateArtistAccountResponse struct {
}

type CreateArtistInviteRequest struct {
	Email     string `json:"email"`
	StageName string `json:"stage_name"`
}

type CreateArtistInviteResponse struct {
	InviteID  string    `json:"invite_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

type RevokeArtistAccountRequest struct {
}
