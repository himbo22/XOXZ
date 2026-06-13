package service

import (
	"context"

	"github.com/himbo22/xoxz/artist-service/internal/logic"
	"github.com/himbo22/xoxz/artist-service/internal/model"
)

type ArtistService interface {
	Create(ctx context.Context, req model.CreateArtistRequest) (*model.ArtistResponse, error)
	GetByID(ctx context.Context, id string) (*model.ArtistResponse, error)
	List(ctx context.Context, req model.ListArtistRequest) (*model.ListArtistResponse, error)
	Update(ctx context.Context, id string, req model.UpdateArtistRequest) (*model.ArtistResponse, error)
	Delete(ctx context.Context, id string) error
}

type artistService struct {
	artistLogic *logic.ArtistLogic
}

func (a *artistService) Create(ctx context.Context, req model.CreateArtistRequest) (*model.ArtistResponse, error) {
	return a.artistLogic.CreateArtist(ctx, req)
}

func (a *artistService) GetByID(ctx context.Context, id string) (*model.ArtistResponse, error) {
	return a.artistLogic.GetArtistByID(ctx, id)
}

func (a *artistService) List(ctx context.Context, req model.ListArtistRequest) (*model.ListArtistResponse, error) {
	return a.artistLogic.ListArtists(ctx, req)
}

func (a *artistService) Update(ctx context.Context, id string, req model.UpdateArtistRequest) (*model.ArtistResponse, error) {
	return a.artistLogic.UpdateArtist(ctx, id, req)
}

func (a *artistService) Delete(ctx context.Context, id string) error {
	return a.artistLogic.DeleteArtist(ctx, id)
}

func NewArtistService(artistLogic *logic.ArtistLogic) ArtistService {
	return &artistService{
		artistLogic: artistLogic,
	}
}
