package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/himbo22/xoxz/account-service/internal/domain/entity"
)

type RoleRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Role, error)
	FindByName(ctx context.Context, name string) (*entity.Role, error)
	Create(ctx context.Context, role *entity.Role) error
	GetAll(ctx context.Context) ([]entity.Role, error)
}
