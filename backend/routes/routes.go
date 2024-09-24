// routes/routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mfuadfakhruzzaki/project/backend/controllers"
	controllers "github.com/mfuadfakhruzzaki/project/backend/middlewares"
	"gorm.io/gorm"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRoutes mengatur semua rute aplikasi
func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	// Public routes
	router.POST("/register", controllers.Register(db))
	router.POST("/login", controllers.Login(db))

	// Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Protected routes
	api := router.Group("/api")
	api.Use(middlewares.AuthMiddleware(db))
	{
		// Admin - User Management
		admin := api.Group("/admin")
		admin.Use(middlewares.AdminMiddleware())
		{
			admin.GET("/users", controllers.GetAllUsers(db))
			admin.DELETE("/users/:id", controllers.DeleteUser(db))
			admin.PUT("/users/:id/status", controllers.UpdateUserStatus(db))
		}

		// Tasks
		tasks := api.Group("/tasks")
		{
			tasks.GET("", controllers.GetTasks(db))
			tasks.POST("", controllers.CreateTask(db))
			tasks.GET("/:id", controllers.GetTaskByID(db))
			tasks.PUT("/:id", controllers.UpdateTask(db))
			tasks.DELETE("/:id", controllers.DeleteTask(db))

			// Assets
			tasks.GET("/:id/assets", controllers.GetAssets(db))
			tasks.POST("/:id/assets", controllers.UploadAsset(db))
		}
	}

	// Other protected routes
	router.GET("/profile", middlewares.AuthMiddleware(db), controllers.GetProfile(db))
	router.PUT("/profile", middlewares.AuthMiddleware(db), controllers.UpdateProfile(db))
	router.GET("/dashboard", middlewares.AuthMiddleware(db), controllers.GetDashboard(db))
}
