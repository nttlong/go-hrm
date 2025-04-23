package db

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"vn.ghrm/internal/config"
	"vn.ghrm/internal/models"
)

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	// Construct the DSN for the hrm database
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPassword, cfg.DbName)

	// Open a connection to PostgreSQL using GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Test the connection
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	log.Println("Successfully connected to PostgreSQL")

	// Run migrations
	if err := MigrateDB(db); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %v", err)
	}

	return db, nil
}

func MigrateDB(db *gorm.DB) error {
	// Auto-migrate the Employee model
	if err := db.AutoMigrate(&models.Employee{}); err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}
	log.Println("Database migration completed")
	return nil
}
