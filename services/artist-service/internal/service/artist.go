package service

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/himbo22/xoxz/artist-service/internal/domain/entity"
	"github.com/himbo22/xoxz/artist-service/internal/domain/repository"
	"github.com/himbo22/xoxz/artist-service/internal/model"
	"github.com/himbo22/xoxz/artist-service/internal/util"
	xoxz "github.com/himbo22/xoxz/common-service/xoxz/logger"
)

type ArtistService interface {
	Create(ctx context.Context, req model.CreateArtistRequest) (*model.ArtistResponse, error)
	GetByID(ctx context.Context, id string) (*model.ArtistResponse, error)
	List(ctx context.Context, req model.ListArtistRequest) (*model.ListArtistResponse, error)
	Update(ctx context.Context, id string, req model.UpdateArtistRequest) (*model.ArtistResponse, error)
	Delete(ctx context.Context, id string) error
}

type artistService struct {
	logger     xoxz.XoxzLogger
	artistRepo repository.ArtistRepository
}

func NewArtistService(
	logger xoxz.XoxzLogger,
	artistRepo repository.ArtistRepository,
) ArtistService {
	return &artistService{
		logger:     logger,
		artistRepo: artistRepo,
	}
}

func (l *artistService) Create(ctx context.Context, req model.CreateArtistRequest) (*model.ArtistResponse, error) {
	stageName := strings.TrimSpace(req.StageName)
	if req.UserID == "" || stageName == "" {
		return nil, util.NewError(http.StatusBadRequest, 4000, "user_id and stage_name are required")
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, util.NewError(http.StatusBadRequest, 4000, "invalid user_id")
	}

	existing, err := l.artistRepo.FindByStageName(ctx, stageName)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, util.NewError(http.StatusConflict, 4090, "stage_name already exists")
	}

	artist := &entity.Artist{
		ID:          uuid.New(),
		UserID:      userID,
		StageName:   stageName,
		DisplayName: req.DisplayName,
		Bio:         req.Bio,
		AvatarURL:   req.AvatarURL,
		BannerURL:   req.BannerURL,
		Status:      "active",
	}

	if err := l.artistRepo.Create(ctx, artist); err != nil {
		return nil, err
	}

	return model.NewArtistResponse(artist), nil
}

func (l *artistService) GetByID(ctx context.Context, id string) (*model.ArtistResponse, error) {
	artistID, err := uuid.Parse(id)
	if err != nil {
		return nil, util.NewError(http.StatusBadRequest, 4000, "invalid artist id")
	}

	artist, err := l.artistRepo.FindByID(ctx, artistID)
	if err != nil {
		return nil, err
	}
	if artist == nil {
		return nil, util.NewError(http.StatusNotFound, 4004, "artist not found")
	}

	return model.NewArtistResponse(artist), nil
}

func (l *artistService) List(ctx context.Context, req model.ListArtistRequest) (*model.ListArtistResponse, error) {
	req.Normalize()
	offset := (req.Page - 1) * req.Limit

	artists, total, err := l.artistRepo.List(ctx, req.Search, req.Limit, offset)
	if err != nil {
		return nil, err
	}

	items := make([]model.ArtistResponse, 0, len(artists))
	for index := range artists {
		items = append(items, *model.NewArtistResponse(&artists[index]))
	}

	return &model.ListArtistResponse{
		Items: items,
		Total: total,
		Page:  req.Page,
		Limit: req.Limit,
	}, nil
}

func (l *artistService) Update(ctx context.Context, id string, req model.UpdateArtistRequest) (*model.ArtistResponse, error) {
	artistID, err := uuid.Parse(id)
	if err != nil {
		return nil, util.NewError(http.StatusBadRequest, 4000, "invalid artist id")
	}

	artist, err := l.artistRepo.FindByID(ctx, artistID)
	if err != nil {
		return nil, err
	}
	if artist == nil {
		return nil, util.NewError(http.StatusNotFound, 4004, "artist not found")
	}

	if req.StageName != nil {
		stageName := strings.TrimSpace(*req.StageName)
		if stageName == "" {
			return nil, util.NewError(http.StatusBadRequest, 4000, "stage_name cannot be empty")
		}
		if stageName != artist.StageName {
			existing, err := l.artistRepo.FindByStageName(ctx, stageName)
			if err != nil {
				return nil, err
			}
			if existing != nil {
				return nil, util.NewError(http.StatusConflict, 4090, "stage_name already exists")
			}
		}
		artist.StageName = stageName
	}
	if req.DisplayName != nil {
		artist.DisplayName = req.DisplayName
	}
	if req.Bio != nil {
		artist.Bio = req.Bio
	}
	if req.AvatarURL != nil {
		artist.AvatarURL = req.AvatarURL
	}
	if req.BannerURL != nil {
		artist.BannerURL = req.BannerURL
	}
	if req.Status != nil {
		artist.Status = *req.Status
	}

	if err := l.artistRepo.Update(ctx, artist); err != nil {
		return nil, err
	}

	return model.NewArtistResponse(artist), nil
}

func (l *artistService) Delete(ctx context.Context, id string) error {
	artistID, err := uuid.Parse(id)
	if err != nil {
		return util.NewError(http.StatusBadRequest, 4000, "invalid artist id")
	}

	artist, err := l.artistRepo.FindByID(ctx, artistID)
	if err != nil {
		return err
	}
	if artist == nil {
		return util.NewError(http.StatusNotFound, 4004, "artist not found")
	}

	return l.artistRepo.Delete(ctx, artistID)
}
