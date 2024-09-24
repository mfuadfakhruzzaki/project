// middlewares/auth.go
package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mfuadfakhruzzaki/project/backend/models"
	"github.com/mfuadfakhruzzaki/project/backend/utils"
	"gorm.io/gorm"
)

// AuthMiddleware verifies the JWT token and sets the current user in context
func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Authorization header missing"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Invalid Authorization header format"})
			c.Abort()
			return
		}

		tokenStr := parts[1]
		userID, err := utils.ParseToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Invalid token"})
			c.Abort()
			return
		}

		var user models.User
		if err := db.Preload("Role").First(&user, userID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "User not found"})
			c.Abort()
			return
		}

		if !user.IsActive {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "User is inactive"})
			c.Abort()
			return
		}

		c.Set("currentUser", user)
		c.Next()
	}
}

// AdminMiddleware ensures that the user has admin privileges
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUser, exists := c.Get("currentUser")
		if !exists {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "User not found"})
			c.Abort()
			return
		}

		user := currentUser.(models.User)
		if user.Role.Name != "admin" {
			c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Admin access required"})
			c.Abort()
			return
		}

		c.Next()
	}
}