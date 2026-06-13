package seeder

import (
	"log"

	"gorm.io/gorm"
)

// Seeder defines the contract every seed file must follow.
type Seeder interface {
	Run(db *gorm.DB) error
	Name() string
}

func StrPtr(s string) *string { return &s }

// RunAll accepts a DB connection and runs all seeders in one transaction.
func RunAll(db *gorm.DB) error {
	// Declare the run order. Tables without foreign keys run first.
	seeders := []Seeder{
		&RoleSeeder{},
		&PermissionSeeder{},
		// Add RolePermissionSeeder here later.
	}

	// Wrap everything in one transaction.
	return db.Transaction(func(tx *gorm.DB) error {
		for _, s := range seeders {
			log.Printf("Running seeder: %s...", s.Name())
			if err := s.Run(tx); err != nil {
				log.Printf("[ERROR] Seeder %s failed: %v", s.Name(), err)
				return err // Returning an error triggers a full rollback.
			}
			log.Printf("Seeder %s completed.", s.Name())
		}
		return nil // Commit transaction
	})
}
