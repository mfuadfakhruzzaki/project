// controllers/tasks.go
package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mfuadfakhruzzaki/project/backend/models"
	"gorm.io/gorm"
)

// GetTasks godoc
// @Summary Mengambil daftar tugas
// @Description Mengambil daftar tugas yang ditugaskan atau dibuat oleh pengguna
// @Tags Tasks
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.Task
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/tasks [get]
func GetTasks(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUserInterface, exists := c.Get("currentUser")
		if !exists {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "User not found"})
			return
		}

		user := currentUserInterface.(models.User)

		var tasks []models.Task
		// Fetch tasks where user is creator or assigned
		if err := db.Preload("Creator").Preload("AssignedTo.User").
			Preload("Comments").
			Preload("Assets").
			Preload("SubTasks").
			Where("created_by = ?", user.ID).
			Or("id IN (SELECT task_id FROM task_assignments WHERE user_id = ?)", user.ID).
			Find(&tasks).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to fetch tasks"})
			return
		}

		c.JSON(http.StatusOK, tasks)
	}
}

// CreateTask godoc
// @Summary Membuat tugas baru
// @Description Membuat tugas baru dalam sistem
// @Tags Tasks
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param task body models.CreateTaskInput true "Create Task"
// @Success 201 {object} models.Task
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/tasks [post]
func CreateTask(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input models.CreateTaskInput
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

		dueDate, err := time.Parse("2006-01-02", input.DueDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid due date format"})
			return
		}

		task := models.Task{
			Title:       input.Title,
			Description: input.Description,
			Priority:    input.Priority,
			Status:      input.Status,
			DueDate:     dueDate,
			CreatedBy:   user.ID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := db.Create(&task).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to create task"})
			return
		}

		// Assign users if any
		for _, userID := range input.AssignedTo {
			// Check if user exists
			var assignedUser models.User
			if err := db.First(&assignedUser, userID).Error; err != nil {
				c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Assigned user not found"})
				return
			}

			assignment := models.TaskAssignment{
				TaskID: task.ID,
				UserID: userID,
			}
			if err := db.Create(&assignment).Error; err != nil {
				c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to assign user to task"})
				return
			}
		}

		// Reload task with associations
		if err := db.Preload("Creator").Preload("AssignedTo.User").Preload("Comments").Preload("Assets").Preload("SubTasks").First(&task, task.ID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to fetch created task"})
			return
		}

		c.JSON(http.StatusCreated, task)
	}
}

// GetTaskByID godoc
// @Summary Mengambil detail tugas
// @Description Mengambil detail tugas berdasarkan ID
// @Tags Tasks
// @Security BearerAuth
// @Param id path int true "Task ID"
// @Produce json
// @Success 200 {object} models.Task
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/tasks/{id} [get]
func GetTaskByID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid task ID"})
			return
		}

		var task models.Task
		if err := db.Preload("Creator").
			Preload("AssignedTo.User").
			Preload("Comments").
			Preload("Assets").
			Preload("SubTasks").
			First(&task, id).Error; err != nil {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Task not found"})
			return
		}

		c.JSON(http.StatusOK, task)
	}
}

// UpdateTask godoc
// @Summary Memperbarui tugas
// @Description Memperbarui informasi tugas berdasarkan ID
// @Tags Tasks
// @Security BearerAuth
// @Param id path int true "Task ID"
// @Param task body models.UpdateTaskInput true "Update Task"
// @Produce json
// @Success 200 {object} models.Task
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/tasks/{id} [put]
func UpdateTask(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid task ID"})
			return
		}

		var input models.UpdateTaskInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}

		var task models.Task
		if err := db.Preload("AssignedTo").First(&task, id).Error; err != nil {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Task not found"})
			return
		}

		// Update fields if provided
		if input.Title != "" {
			task.Title = input.Title
		}
		if input.Description != "" {
			task.Description = input.Description
		}
		if input.Priority != "" {
			task.Priority = input.Priority
		}
		if input.Status != "" {
			task.Status = input.Status
		}
		if input.DueDate != "" {
			dueDate, err := time.Parse("2006-01-02", input.DueDate)
			if err != nil {
				c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid due date format"})
				return
			}
			task.DueDate = dueDate
		}

		task.UpdatedAt = time.Now()

		if err := db.Save(&task).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to update task"})
			return
		}

		// Update assignments if provided
		if input.AssignedTo != nil {
			// Delete existing assignments
			db.Where("task_id = ?", task.ID).Delete(&models.TaskAssignment{})

			// Assign new users
			for _, userID := range input.AssignedTo {
				// Check if user exists
				var assignedUser models.User
				if err := db.First(&assignedUser, userID).Error; err != nil {
					c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Assigned user not found"})
					return
				}

				assignment := models.TaskAssignment{
					TaskID: task.ID,
					UserID: userID,
				}
				if err := db.Create(&assignment).Error; err != nil {
					c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to assign user to task"})
					return
				}
			}
		}

		// Reload task with associations
		if err := db.Preload("Creator").Preload("AssignedTo.User").Preload("Comments").Preload("Assets").Preload("SubTasks").First(&task, task.ID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to fetch updated task"})
			return
		}

		c.JSON(http.StatusOK, task)
	}
}

// DeleteTask godoc
// @Summary Menghapus tugas
// @Description Menghapus tugas berdasarkan ID
// @Tags Tasks
// @Security BearerAuth
// @Param id path int true "Task ID"
// @Produce json
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/tasks/{id} [delete]
func DeleteTask(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid task ID"})
			return
		}

		var task models.Task
		if err := db.First(&task, id).Error; err != nil {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Task not found"})
			return
		}

		if err := db.Delete(&task).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to delete task"})
			return
		}

		c.JSON(http.StatusOK, models.SuccessResponse{Message: "Task deleted successfully"})
	}
}
