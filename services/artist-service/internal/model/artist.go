package model

import (
	"time"

	"github.com/google/uuid"
	_const "github.com/himbo22/xoxz/artist-service/internal/const"
	"github.com/himbo22/xoxz/artist-service/internal/domain/entity"
)

type CreateArtistRequest struct {
	UserID      string  `json:"user_id"`
	StageName   string  `json:"stage_name"`
	DisplayName *string `json:"display_name"`
	Bio         *string `json:"bio"`
	AvatarURL   *string `json:"avatar_url"`
	BannerURL   *string `json:"banner_url"`
}

type UpdateArtistRequest struct {
	StageName   *string `json:"stage_name"`
	DisplayName *string `json:"display_name"`
	Bio         *string `json:"bio"`
	AvatarURL   *string `json:"avatar_url"`
	BannerURL   *string `json:"banner_url"`
	Status      *string `json:"status"`
}

type ListArtistRequest struct {
	Page   int
	Limit  int
	Search string
}

func (r *ListArtistRequest) Normalize() {
	if r.Page <= 0 {
		r.Page = _const.DefaultPage
	}
	if r.Limit <= 0 {
		r.Limit = _const.DefaultLimit
	}
	if r.Limit > _const.MaxLimit {
		r.Limit = _const.MaxLimit
	}
}

type ArtistResponse struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	StageName   string    `json:"stage_name"`
	DisplayName *string   `json:"display_name"`
	Bio         *string   `json:"bio"`
	AvatarURL   *string   `json:"avatar_url"`
	BannerURL   *string   `json:"banner_url"`
	Verified    bool      `json:"verified"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ListArtistResponse struct {
	Items []ArtistResponse `json:"items"`
	Total int64            `json:"total"`
	Page  int              `json:"page"`
	Limit int              `json:"limit"`
}

func NewArtistResponse(artist *entity.Artist) *ArtistResponse {
	return &ArtistResponse{
		ID:          artist.ID,
		UserID:      artist.UserID,
		StageName:   artist.StageName,
		DisplayName: artist.DisplayName,
		Bio:         artist.Bio,
		AvatarURL:   artist.AvatarURL,
		BannerURL:   artist.BannerURL,
		Verified:    artist.Verified,
		Status:      artist.Status,
		CreatedAt:   artist.CreatedAt,
		UpdatedAt:   artist.UpdatedAt,
	}
}
