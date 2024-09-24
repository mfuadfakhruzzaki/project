// controllers/auth.go
package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mfuadfakhruzzaki/project/backend/models"
	"github.com/mfuadfakhruzzaki/project/backend/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Register godoc
// @Summary Registrasi pengguna baru
// @Description Mendaftarkan pengguna baru ke sistem
// @Tags Auth
// @Accept json
// @Produce json
// @Param register body models.RegisterInput true "Register Input"
// @Success 201 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /register [post]
func Register(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input models.RegisterInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}

		// Check if email already exists
		var existingUser models.User
		if err := db.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Email already registered"})
			return
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to hash password"})
			return
		}

		// Assign default role (assume role_id = 1 adalah 'user')
		var role models.Role
		if err := db.First(&role, 1).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Default role not found"})
			return
		}

		user := models.User{
			Username:  input.Username,
			Email:     input.Email,
			Password:  string(hashedPassword),
			IsActive:  true,
			RoleID:    role.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to create user"})
			return
		}

		c.JSON(http.StatusCreated, models.SuccessResponse{Message: "User registered successfully"})
	}
}

// Login godoc
// @Summary Login pengguna
// @Description Mengautentikasi pengguna dan menghasilkan token JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param login body models.LoginInput true "Login Input"
// @Success 200 {object} models.TokenResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /login [post]
func Login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input models.LoginInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}

		var user models.User
		if err := db.Preload("Role").Where("email = ?", input.Email).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Invalid email or password"})
			return
		}

		// Compare password
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Invalid email or password"})
			return
		}

		if !user.IsActive {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "User is inactive"})
			return
		}

		// Generate JWT token
		token, err := utils.GenerateToken(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, models.TokenResponse{Token: token})
	}
}
