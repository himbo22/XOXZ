package repository_impl

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/himbo22/xoxz/account-service/internal/domain/entity"
	"github.com/himbo22/xoxz/account-service/internal/domain/repository"
	"gorm.io/gorm"
)

type permissionRepository struct {
	db *gorm.DB
}

func (r *permissionRepository) GetAll(ctx context.Context) ([]entity.Permission, error) {
	var rolePermissions []entity.Permission

	err := r.db.WithContext(ctx).
		Find(&rolePermissions).
		Error

	return rolePermissions, err
}

func NewPermissionRepository(db *gorm.DB) repository.PermissionRepository {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Permission, error) {
	var permission entity.Permission
	err := GetDB(ctx, r.db).Where("id = ?", id).First(&permission).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *permissionRepository) FindByCode(ctx context.Context, code string) (*entity.Permission, error) {
	var permission entity.Permission
	err := GetDB(ctx, r.db).Where("code = ?", code).First(&permission).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *permissionRepository) Create(ctx context.Context, permission *entity.Permission) error {
	return GetDB(ctx, r.db).Create(permission).Error
}

func (r *permissionRepository) Update(ctx context.Context, permission *entity.Permission) error {
	return GetDB(ctx, r.db).Save(permission).Error
}

func (r *permissionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return GetDB(ctx, r.db).Delete(&entity.Permission{}, "id = ?", id).Error
}
