package seeder

import (
	"github.com/himbo22/xoxz/account-service/internal/domain/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PermissionSeeder struct{}

func (s *PermissionSeeder) Name() string {
	return "PermissionSeeder"
}

func (s *PermissionSeeder) Run(db *gorm.DB) error {
	permissions := []entity.Permission{
		{Code: "TEST1", Description: StrPtr("McLaren")},
		{Code: "TEST2", Description: StrPtr("Ferrari")},
	}

	// UPSERT: if the code conflicts, update the description.
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "code"}},
		DoUpdates: clause.AssignmentColumns([]string{"description"}),
	}).Create(&permissions).Error
}
