package repository_impl

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/himbo22/xoxz/account-service/internal/domain/entity"
	"github.com/himbo22/xoxz/account-service/internal/domain/repository"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{db: db}
}

// UpdateLastLogin implements repository.UserRepository.
func (r *userRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	return GetDB(ctx, r.db).Model(&entity.User{}).
		Where("id = ?", id).
		Update("last_login_at", gorm.Expr("NOW()")).Error
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := GetDB(ctx, r.db).Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := GetDB(ctx, r.db).Where("id = ?", id).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	err := GetDB(ctx, r.db).Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	return GetDB(ctx, r.db).Create(user).Error
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	return GetDB(ctx, r.db).Save(user).Error
}
