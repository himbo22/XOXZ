package repository_impl

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/himbo22/xoxz/account-service/internal/domain/entity"
	"github.com/himbo22/xoxz/account-service/internal/domain/repository"
	"gorm.io/gorm"
)

type roleRepository struct {
	db *gorm.DB
}

func (r *roleRepository) GetAll(ctx context.Context) ([]entity.Role, error) {
	//TODO implement me
	panic("implement me")
}

func NewRoleRepository(db *gorm.DB) repository.RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Role, error) {
	var role entity.Role
	err := GetDB(ctx, r.db).Where("id = ?", id).First(&role).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) FindByName(ctx context.Context, name string) (*entity.Role, error) {
	var role entity.Role
	err := GetDB(ctx, r.db).Where("name = ?", name).First(&role).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) Create(ctx context.Context, role *entity.Role) error {
	return GetDB(ctx, r.db).Create(role).Error
}
