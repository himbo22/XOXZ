package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/himbo22/xoxz/account-service/internal/domain/entity"
)

type UserRoleRepository interface {
	FindByID(ctx context.Context, roleID, userID uuid.UUID) (*entity.UserRole, error)
	Create(ctx context.Context, userRole *entity.UserRole) error
	Update(ctx context.Context, userRole *entity.UserRole) error
	Delete(ctx context.Context, roleID, userID uuid.UUID) error
}
