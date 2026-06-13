package repository_impl

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/himbo22/xoxz/account-service/internal/domain/entity"
	"github.com/himbo22/xoxz/account-service/internal/domain/repository"
	"gorm.io/gorm"
)

type userRoleRepository struct {
	db *gorm.DB
}

func NewUserRoleRepository(db *gorm.DB) repository.UserRoleRepository {
	return &userRoleRepository{db: db}
}

func (r *userRoleRepository) FindByID(ctx context.Context, roleID, userID uuid.UUID) (*entity.UserRole, error) {
	var userRole entity.UserRole
	err := GetDB(ctx, r.db).
		Where("role_id = ? AND user_id = ?", roleID, userID).
		First(&userRole).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &userRole, nil
}

func (r *userRoleRepository) Create(ctx context.Context, userRole *entity.UserRole) error {
	return GetDB(ctx, r.db).Create(userRole).Error
}

func (r *userRoleRepository) Update(ctx context.Context, userRole *entity.UserRole) error {
	return GetDB(ctx, r.db).Save(userRole).Error
}

func (r *userRoleRepository) Delete(ctx context.Context, roleID, userID uuid.UUID) error {
	return GetDB(ctx, r.db).
		Delete(&entity.UserRole{}, "role_id = ? AND user_id = ?", roleID, userID).Error
}
