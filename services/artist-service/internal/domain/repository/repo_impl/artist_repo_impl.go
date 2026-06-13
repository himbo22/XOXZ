package repo_impl

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/himbo22/xoxz/artist-service/internal/domain/entity"
	"github.com/himbo22/xoxz/artist-service/internal/domain/repository"
	"gorm.io/gorm"
)

type artistRepository struct {
	db *gorm.DB
}

func NewArtistRepository(db *gorm.DB) repository.ArtistRepository {
	return &artistRepository{db: db}
}

func (r *artistRepository) Create(ctx context.Context, artist *entity.Artist) error {
	return r.db.WithContext(ctx).Create(artist).Error
}

func (r *artistRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Artist, error) {
	var artist entity.Artist
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&artist).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &artist, nil
}

func (r *artistRepository) FindByStageName(ctx context.Context, stageName string) (*entity.Artist, error) {
	var artist entity.Artist
	err := r.db.WithContext(ctx).Where("stage_name = ?", stageName).First(&artist).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &artist, nil
}

func (r *artistRepository) List(ctx context.Context, search string, limit, offset int) ([]entity.Artist, int64, error) {
	var artists []entity.Artist
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Artist{})
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("stage_name ILIKE ? OR display_name ILIKE ?", like, like)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&artists).Error; err != nil {
		return nil, 0, err
	}

	return artists, total, nil
}

func (r *artistRepository) Update(ctx context.Context, artist *entity.Artist) error {
	return r.db.WithContext(ctx).Save(artist).Error
}

func (r *artistRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Artist{}, "id = ?", id).Error
}
