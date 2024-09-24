// config/config.go
package config

import (
	"log"
	"os"

	"github.com/mfuadfakhruzzaki/project/backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SetupDatabase initializes the database connection
func SetupDatabase() *gorm.DB {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set in environment")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	return db
}

// MigrateDatabase runs the database migrations
func MigrateDatabase(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.Role{},
		&models.User{},
		&models.Task{},
		&models.TaskAssignment{},
		&models.Asset{},
		&models.Comment{},
		&models.SubTask{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}
