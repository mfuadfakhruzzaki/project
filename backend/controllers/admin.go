package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mfuadfakhruzzaki/project/backend/models"
	"gorm.io/gorm"
)

// GetAllUsers godoc
// @Summary Mengambil daftar semua pengguna
// @Description Mengambil daftar semua pengguna dengan peran mereka
// @Tags Admin - User Management
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.User
// @Failure 500 {object} models.ErrorResponse
// @Router /api/admin/users [get]
func GetAllUsers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []models.User
		if err := db.Preload("Role").Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to fetch users"})
			return
		}
		c.JSON(http.StatusOK, users)
	}
}

// DeleteUser godoc
// @Summary Menghapus pengguna berdasarkan ID
// @Description Menghapus pengguna berdasarkan ID
// @Tags Admin - User Management
// @Security BearerAuth
// @Param id path int true "User ID"
// @Produce json
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/admin/users/{id} [delete]
func DeleteUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid user ID"})
			return
		}

		// Prevent deleting admin user
		var user models.User
		if err := db.Preload("Role").First(&user, id).Error; err != nil {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "User not found"})
			return
		}

		if user.Role.Name == "admin" {
			c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Cannot delete admin user"})
			return
		}

		if err := db.Delete(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to delete user"})
			return
		}

		c.JSON(http.StatusOK, models.SuccessResponse{Message: "User deleted successfully"})
	}
}

// UpdateUserStatus godoc
// @Summary Memperbarui status aktif pengguna berdasarkan ID
// @Description Memperbarui status aktif (is_active) pengguna berdasarkan ID
// @Tags Admin - User Management
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param status body models.UpdateUserStatusRequest true "Status Update"
// @Produce json
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/admin/users/{id}/status [put]
func UpdateUserStatus(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid user ID"})
			return
		}

		var input models.UpdateUserStatusRequest
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}

		var user models.User
		if err := db.First(&user, id).Error; err != nil {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "User not found"})
			return
		}

		user.IsActive = input.IsActive
		if err := db.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to update user status"})
			return
		}

		c.JSON(http.StatusOK, models.SuccessResponse{Message: "User status updated successfully"})
	}
}
