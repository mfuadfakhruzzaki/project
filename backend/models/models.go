// models/models.go
package models

import (
	"time"
)

// Role represents the user roles in the system
type Role struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name" binding:"required"`
}

// User represents a user in the system
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" binding:"required"`
	Email     string    `json:"email" gorm:"unique" binding:"required,email"`
	Password  string    `json:"-" binding:"required,min=6"`
	IsActive  bool      `json:"is_active"`
	RoleID    uint      `json:"role_id"`
	Role      Role      `json:"role" gorm:"foreignKey:RoleID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Task represents a task in the system
type Task struct {
	ID          uint             `json:"id" gorm:"primaryKey"`
	Title       string           `json:"title" binding:"required"`
	Description string           `json:"description"`
	Priority    string           `json:"priority" binding:"oneof=high medium normal low"`
	Status      string           `json:"status" binding:"oneof=todo in_progress completed"`
	DueDate     time.Time        `json:"due_date" binding:"required"`
	CreatedBy   uint             `json:"created_by"`
	Creator     User             `json:"creator" gorm:"foreignKey:CreatedBy"`
	AssignedTo  []TaskAssignment `json:"assigned_to" gorm:"foreignKey:TaskID"`
	Comments    []Comment        `json:"comments" gorm:"foreignKey:TaskID"`
	Assets      []Asset          `json:"assets" gorm:"foreignKey:TaskID"`
	SubTasks    []SubTask        `json:"sub_tasks" gorm:"foreignKey:TaskID"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

// TaskAssignment represents the assignment of a task to a user
type TaskAssignment struct {
	ID     uint `json:"id" gorm:"primaryKey"`
	TaskID uint `json:"task_id"`
	UserID uint `json:"user_id"`
	User   User `json:"user" gorm:"foreignKey:UserID"`
}

// Asset represents an asset associated with a task
type Asset struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	FilePath   string    `json:"file_path"`
	TaskID     uint      `json:"task_id"`
	UploadedBy uint      `json:"uploaded_by"`
	UploadedAt time.Time `json:"uploaded_at"`
}

// Comment represents a comment on a task
type Comment struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Content   string    `json:"content" binding:"required"`
	TaskID    uint      `json:"task_id"`
	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SubTask represents a sub-task within a main task
type SubTask struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	Priority    string    `json:"priority" binding:"oneof=high medium normal low"`
	Status      string    `json:"status" binding:"oneof=todo in_progress completed"`
	DueDate     time.Time `json:"due_date" binding:"required"`
	TaskID      uint      `json:"task_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Response Models

// ErrorResponse represents the structure of error responses
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse represents the structure of success responses
type SuccessResponse struct {
	Message string `json:"message"`
}

// TokenResponse represents the structure of JWT token responses
type TokenResponse struct {
	Token string `json:"token"`
}

// RegisterInput represents the input for user registration
type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginInput represents the input for user login
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UpdateProfileInput represents the input for updating user profile
type UpdateProfileInput struct {
	Username string `json:"username"`
	Email    string `json:"email" binding:"email"`
	Password string `json:"password" binding:"min=6"`
}

// UpdateTaskInput represents the input for updating a task
type UpdateTaskInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority" binding:"oneof=high medium normal low"`
	Status      string `json:"status" binding:"oneof=todo in_progress completed"`
	DueDate     string `json:"due_date" binding:"required"`
	AssignedTo  []uint `json:"assigned_to"`
}

// UpdateUserStatusRequest represents the input for updating user status
type UpdateUserStatusRequest struct {
	IsActive bool `json:"is_active" binding:"required"`
}

// CreateTaskInput represents the input for creating a new task
type CreateTaskInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Priority    string `json:"priority" binding:"oneof=high medium normal low"`
	Status      string `json:"status" binding:"oneof=todo in_progress completed"`
	DueDate     string `json:"due_date" binding:"required"`
	AssignedTo  []uint `json:"assigned_to"`
}

// DashboardResponse represents the response structure for dashboard data
type DashboardResponse struct {
	Todo       int `json:"todo"`
	InProgress int `json:"in_progress"`
	Completed  int `json:"completed"`
}

// AssetResponse represents the response structure for assets
type AssetResponse struct {
	ID       uint   `json:"id"`
	FilePath string `json:"file_path"`
}

// TaskResponse represents the response structure for a task
type TaskResponse struct {
	Task
}

// UserResponse represents the response structure for a user
type UserResponse struct {
	User
}
