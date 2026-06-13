package repository

import (
	"context"

	"github.com/himbo22/xoxz/livestream-service/internal/domain/entity"
	"github.com/himbo22/xoxz/livestream-service/internal/model"
)

type LivestreamRepository interface {
	Create(ctx context.Context, s *entity.LivestreamRoom) error
	GetByID(ctx context.Context, id string) (*entity.LivestreamRoom, error)
	UpdateStatus(ctx context.Context, id string, status model.StreamStatus) error
}
