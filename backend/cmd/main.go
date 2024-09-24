// cmd/main.go
package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/mfuadfakhruzzaki/project/backend/config"
	"github.com/mfuadfakhruzzaki/project/backend/routes"

	// Swagger docs
	_ "github.com/mfuadfakhruzzaki/project/backend/docs"
)

// @title Project Management API
// @version 1.0
// @description API untuk mengelola proyek dan tugas.
// @termsOfService http://example.com/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	// Initialize Database
	db := config.SetupDatabase()
	// Migrate models
	config.MigrateDatabase(db)
	// Seed roles
	config.SeedRoles(db)

	// Set Gin to release mode in production
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Initialize Routes
	routes.SetupRoutes(router, db)

	// Run the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
