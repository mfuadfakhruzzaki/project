// controllers/assets.go
package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mfuadfakhruzzaki/project/backend/models"
	"github.com/mfuadfakhruzzaki/project/backend/utils"
	"gorm.io/gorm"
)

// GetAssets godoc
// @Summary Mengambil semua aset terkait tugas
// @Description Mengambil semua aset yang terkait dengan tugas berdasarkan ID tugas
// @Tags Assets
// @Security BearerAuth
// @Param id path int true "Task ID"
// @Produce json
// @Success 200 {array} models.AssetResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/tasks/{id}/assets [get]
func GetAssets(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		taskIDParam := c.Param("id")
		taskID, err := strconv.Atoi(taskIDParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid task ID"})
			return
		}

		var task models.Task
		if err := db.First(&task, taskID).Error; err != nil {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Task not found"})
			return
		}

		var assets []models.Asset
		if err := db.Where("task_id = ?", taskID).Find(&assets).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to fetch assets"})
			return
		}

		var assetResponses []models.AssetResponse
		for _, asset := range assets {
			assetResponses = append(assetResponses, models.AssetResponse{
				ID:       asset.ID,
				FilePath: asset.FilePath,
			})
		}

		c.JSON(http.StatusOK, assetResponses)
	}
}

// UploadAsset godoc
// @Summary Mengunggah aset ke tugas
// @Description Mengunggah file aset dan mengaitkannya dengan tugas berdasarkan ID tugas
// @Tags Assets
// @Security BearerAuth
// @Param id path int true "Task ID"
// @Param file formData file true "File to upload"
// @Produce json
// @Success 201 {object} models.AssetResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/tasks/{id}/assets [post]
func UploadAsset(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		taskIDParam := c.Param("id")
		taskID, err := strconv.Atoi(taskIDParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid task ID"})
			return
		}

		var task models.Task
		if err := db.First(&task, taskID).Error; err != nil {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Task not found"})
			return
		}

		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "File is required"})
			return
		}

		// Upload file to a directory (e.g., ./uploads)
		uploadPath := "./uploads/"
		if err := utils.CreateDirIfNotExists(uploadPath); err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to create upload directory"})
			return
		}

		fileName := strconv.Itoa(int(task.ID)) + "_" + file.Filename
		fullPath := uploadPath + fileName

		if err := c.SaveUploadedFile(file, fullPath); err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to save file"})
			return
		}

		// Get current user
		currentUserInterface, exists := c.Get("currentUser")
		if !exists {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "User not found"})
			return
		}

		user := currentUserInterface.(models.User)

		// Create Asset record
		asset := models.Asset{
			FilePath:   fullPath,
			TaskID:     uint(taskID),
			UploadedBy: user.ID,
			UploadedAt: time.Now(),
		}

		if err := db.Create(&asset).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to save asset"})
			return
		}

		c.JSON(http.StatusCreated, models.AssetResponse{
			ID:       asset.ID,
			FilePath: asset.FilePath,
		})
	}
}
