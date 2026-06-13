package repository_impl

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/himbo22/xoxz/account-service/internal/domain/entity"
	"github.com/himbo22/xoxz/account-service/internal/domain/repository"
	"gorm.io/gorm"
)

type rolePermissionRepository struct {
	db *gorm.DB
}

func (r *rolePermissionRepository) GetAll(ctx context.Context) ([]entity.RolePermission, error) {
	var rolePermissions []entity.RolePermission

	err := r.db.WithContext(ctx).
		Find(&rolePermissions).
		Error

	return rolePermissions, err
}

func NewRolePermissionRepository(db *gorm.DB) repository.RolePermissionRepository {
	return &rolePermissionRepository{db: db}
}

func (r *rolePermissionRepository) FindByID(ctx context.Context, roleID, permissionID uuid.UUID) (*entity.RolePermission, error) {
	var rolePermission entity.RolePermission
	err := GetDB(ctx, r.db).
		Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		First(&rolePermission).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &rolePermission, nil
}

func (r *rolePermissionRepository) Create(ctx context.Context, rolePermission *entity.RolePermission) error {
	return GetDB(ctx, r.db).Create(rolePermission).Error
}

func (r *rolePermissionRepository) Update(ctx context.Context, rolePermission *entity.RolePermission) error {
	return GetDB(ctx, r.db).Save(rolePermission).Error
}

func (r *rolePermissionRepository) Delete(ctx context.Context, roleID, permissionID uuid.UUID) error {
	return GetDB(ctx, r.db).
		Delete(&entity.RolePermission{}, "role_id = ? AND permission_id = ?", roleID, permissionID).Error
}
