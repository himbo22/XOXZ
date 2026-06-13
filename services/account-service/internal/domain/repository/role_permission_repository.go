package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/himbo22/xoxz/account-service/internal/domain/entity"
)

type RolePermissionRepository interface {
	FindByID(ctx context.Context, roleID, permissionID uuid.UUID) (*entity.RolePermission, error)
	Create(ctx context.Context, rolePermission *entity.RolePermission) error
	Update(ctx context.Context, rolePermission *entity.RolePermission) error
	Delete(ctx context.Context, roleID, permissionID uuid.UUID) error
	GetAll(ctx context.Context) ([]entity.RolePermission, error)
}
