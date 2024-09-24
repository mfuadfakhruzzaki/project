// controllers/dashboard.go
package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mfuadfakhruzzaki/project/backend/models"
	"gorm.io/gorm"
)

// GetDashboard godoc
// @Summary Mengambil data dashboard pengguna
// @Description Mengambil jumlah tugas berdasarkan status (todo, in_progress, completed) untuk pengguna yang sedang login
// @Tags Dashboard
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.DashboardResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /dashboard [get]
func GetDashboard(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUserInterface, exists := c.Get("currentUser")
		if !exists {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "User not found"})
			return
		}

		user := currentUserInterface.(models.User)

		var dashboard models.DashboardResponse

		// Count todo
		if err := db.Model(&models.Task{}).
			Where("created_by = ? OR id IN (SELECT task_id FROM task_assignments WHERE user_id = ?)", user.ID, user.ID).
			Where("status = ?", "todo").
			Count(&dashboard.Todo).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to count todo tasks"})
			return
		}

		// Count in_progress
		if err := db.Model(&models.Task{}).
			Where("created_by = ? OR id IN (SELECT task_id FROM task_assignments WHERE user_id = ?)", user.ID, user.ID).
			Where("status = ?", "in_progress").
			Count(&dashboard.InProgress).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to count in_progress tasks"})
			return
		}

		// Count completed
		if err := db.Model(&models.Task{}).
			Where("created_by = ? OR id IN (SELECT task_id FROM task_assignments WHERE user_id = ?)", user.ID, user.ID).
			Where("status = ?", "completed").
			Count(&dashboard.Completed).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to count completed tasks"})
			return
		}

		c.JSON(http.StatusOK, dashboard)
	}
}
