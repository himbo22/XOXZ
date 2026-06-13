package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/himbo22/xoxz/artist-service/internal/domain/entity"
)

type ArtistRepository interface {
	Create(ctx context.Context, artist *entity.Artist) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Artist, error)
	FindByStageName(ctx context.Context, stageName string) (*entity.Artist, error)
	List(ctx context.Context, search string, limit, offset int) ([]entity.Artist, int64, error)
	Update(ctx context.Context, artist *entity.Artist) error
	Delete(ctx context.Context, id uuid.UUID) error
}
