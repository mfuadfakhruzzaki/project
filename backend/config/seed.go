// config/seed.go
package config

import (
	"log"

	"github.com/mfuadfakhruzzaki/project/backend/models"
	"gorm.io/gorm"
)

// SeedRoles seeds the roles into the database
func SeedRoles(db *gorm.DB) {
	var count int64
	db.Model(&models.Role{}).Count(&count)
	if count == 0 {
		roles := []models.Role{
			{ID: 1, Name: "user"},
			{ID: 2, Name: "admin"},
		}
		if err := db.Create(&roles).Error; err != nil {
			log.Fatalf("Failed to seed roles: %v", err)
		}
		log.Println("Seeded roles successfully")
	}
}
