// controllers/profile.go
package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mfuadfakhruzzaki/project/backend/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// GetProfile godoc
// @Summary Mengambil profil pengguna
// @Description Mengambil profil pengguna yang sedang login
// @Tags Profile
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.User
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /profile [get]
func GetProfile(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUserInterface, exists := c.Get("currentUser")
		if !exists {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "User not found"})
			return
		}

		user := currentUserInterface.(models.User)

		var fullUser models.User
		if err := db.Preload("Role").First(&fullUser, user.ID).Error; err != nil {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "User not found"})
			return
		}

		c.JSON(http.StatusOK, fullUser)
	}
}

// UpdateProfile godoc
// @Summary Memperbarui profil pengguna
// @Description Memperbarui informasi profil pengguna yang sedang login
// @Tags Profile
// @Security BearerAuth
// @Param profile body models.UpdateProfileInput true "Update Profile"
// @Produce json
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /profile [put]
func UpdateProfile(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input models.UpdateProfileInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}

		currentUserInterface, exists := c.Get("currentUser")
		if !exists {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "User not found"})
			return
		}

		user := currentUserInterface.(models.User)

		// Update fields if provided
		if input.Username != "" {
			user.Username = input.Username
		}
		if input.Email != "" {
			// Check if email is taken
			var existingUser models.User
			if err := db.Where("email = ? AND id != ?", input.Email, user.ID).First(&existingUser).Error; err == nil {
				c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Email already in use"})
				return
			}
			user.Email = input.Email
		}
		if input.Password != "" {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to hash password"})
				return
			}
			user.Password = string(hashedPassword)
		}

		user.UpdatedAt = time.Now()

		if err := db.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to update profile"})
			return
		}

		c.JSON(http.StatusOK, models.SuccessResponse{Message: "Profile updated successfully"})
	}
}
