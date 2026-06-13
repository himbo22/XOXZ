package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/himbo22/xoxz/account-service/internal/domain/entity"
)

type PermissionRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Permission, error)
	FindByCode(ctx context.Context, code string) (*entity.Permission, error)
	Create(ctx context.Context, permission *entity.Permission) error
	Update(ctx context.Context, permission *entity.Permission) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetAll(ctx context.Context) ([]entity.Permission, error)
}
