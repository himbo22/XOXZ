package seeder

import (
	_const "github.com/himbo22/xoxz/account-service/internal/const"
	"github.com/himbo22/xoxz/account-service/internal/domain/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RoleSeeder struct{}

func (s *RoleSeeder) Name() string {
	return "RoleSeeder"
}

func (s *RoleSeeder) Run(db *gorm.DB) error {
	roles := []entity.Role{
		{Name: _const.ROLE_ADMIN, Description: StrPtr("Top-level administrator")},
		{Name: _const.ROLE_USER, Description: StrPtr("Standard user")},
		{Name: _const.ROLE_IDOL, Description: StrPtr("Idol Kpop")},
	}

	// UPSERT: if the name conflicts, update the description.
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoUpdates: clause.AssignmentColumns([]string{"description"}),
	}).Create(&roles).Error
}
