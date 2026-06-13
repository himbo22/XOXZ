package repository_impl

import (
	"context"
	"errors"

	"github.com/himbo22/xoxz/account-service/internal/domain/entity"
	"github.com/himbo22/xoxz/account-service/internal/domain/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type identityRepository struct {
	db *gorm.DB
}

func NewIdentityRepository(db *gorm.DB) repository.IdentityRepository {
	return &identityRepository{db: db}
}

// FindBySubjectAndProvider implements repository.IdentityRepository.
func (i *identityRepository) FindBySubjectAndProvider(ctx context.Context, subject string, provider string) (*entity.Identity, error) {
	identity := entity.Identity{}
	err := GetDB(ctx, i.db).Where("provider_user_id = ? AND provider = ?", subject, provider).First(&identity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &identity, nil
}

// Upsert implements repository.IdentityRepository.
func (i *identityRepository) Upsert(ctx context.Context, identity *entity.Identity) error {
	return GetDB(ctx, i.db).Clauses(clause.OnConflict{
		// 1. Specify the conflicting columns, which require a unique index in the database.
		Columns: []clause.Column{{Name: "provider"}, {Name: "provider_user_id"}},
		TargetWhere: clause.Where{
			Exprs: []clause.Expression{
				clause.Expr{SQL: "deleted_at IS NULL"},
			},
		},
		// 2. On conflict, update the JSON payload and timestamp.
		DoUpdates: clause.AssignmentColumns([]string{"provider_data", "updated_at"}),
	}).Create(identity).Error
}
