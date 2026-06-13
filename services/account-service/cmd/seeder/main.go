package main

import (
	"fmt"
	"log"

	"github.com/himbo22/xoxz/account-service/internal/config"
	"github.com/himbo22/xoxz/account-service/internal/domain/seeder"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	loadConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	dbConfig := loadConfig.Database
	// 1. Initialize the connection, usually from config or environment.
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=%v", dbConfig.Host, dbConfig.User, dbConfig.Password, dbConfig.DBName, dbConfig.Port, dbConfig.SSLMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	log.Println("Starting core data seeding...")

	// 2. Run all seeders.
	if err := seeder.RunAll(db); err != nil {
		log.Fatalf("Seeding failed. The database transaction was rolled back. Error: %v", err)
	}

	log.Println("Seeding completed successfully!")
}
