package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/anisetusbambangmanalu19/moonlay-task-management-app/backend/internal/models"
	"github.com/anisetusbambangmanalu19/moonlay-task-management-app/backend/internal/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TaskHandler handles task-related HTTP requests
type TaskHandler struct {
	taskRepo *repository.TaskRepository
	userRepo *repository.UserRepository
}

// NewTaskHandler creates a new TaskHandler
func NewTaskHandler(taskRepo *repository.TaskRepository, userRepo *repository.UserRepository) *TaskHandler {
	return &TaskHandler{taskRepo: taskRepo, userRepo: userRepo}
}

// --- Request structs ---

type CreateTaskRequest struct {
	Title       string    `json:"title" binding:"required,min=1,max=200"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Deadline    time.Time `json:"deadline" binding:"required"`
	AssigneeID  int64     `json:"assignee_id" binding:"required,min=1"`
}

type UpdateTaskRequest struct {
	Title       string    `json:"title" binding:"required,min=1,max=200"`
	Description string    `json:"description"`
	Status      string    `json:"status" binding:"required"`
	Deadline    time.Time `json:"deadline" binding:"required"`
	AssigneeID  int64     `json:"assignee_id" binding:"required,min=1"`
}

type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// --- Handlers ---

// GetAll returns all tasks with assignee info.
// GET /api/tasks
func (h *TaskHandler) GetAll(c *gin.Context) {
	tasks, err := h.taskRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data task"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tasks})
}

// GetByID returns a single task by ID.
// GET /api/tasks/:id
func (h *TaskHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID task tidak valid"})
		return
	}

	task, err := h.taskRepo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": task})
}

// Create creates a new task.
// POST /api/tasks
func (h *TaskHandler) Create(c *gin.Context) {
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate assignee exists
	if _, err := h.userRepo.FindByID(req.AssigneeID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Assignee tidak ditemukan"})
		return
	}

	// Set default status and validate
	status := models.TaskStatus(req.Status)
	if req.Status == "" {
		status = models.TaskStatusTodo
	}
	if !status.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status tidak valid. Pilihan: todo, in_progress, done"})
		return
	}

	// Get current user from JWT context
	userID := c.GetInt64("user_id")

	task := &models.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      status,
		Deadline:    req.Deadline,
		AssigneeID:  req.AssigneeID,
		CreatedBy:   &userID,
	}

	if err := h.taskRepo.Create(task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat task"})
		return
	}

	// Reload with associations for response
	created, _ := h.taskRepo.FindByID(task.ID)
	c.JSON(http.StatusCreated, gin.H{
		"message": "Task berhasil dibuat",
		"data":    created,
	})
}

// Update updates all fields of a task.
// PUT /api/tasks/:id
func (h *TaskHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID task tidak valid"})
		return
	}

	task, err := h.taskRepo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil task"})
		return
	}

	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate status
	status := models.TaskStatus(req.Status)
	if !status.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status tidak valid. Pilihan: todo, in_progress, done"})
		return
	}

	// Validate assignee
	if _, err := h.userRepo.FindByID(req.AssigneeID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Assignee tidak ditemukan"})
		return
	}

	task.Title = req.Title
	task.Description = req.Description
	task.Status = status
	task.Deadline = req.Deadline
	task.AssigneeID = req.AssigneeID
	task.Assignee = nil // clear to prevent GORM association update issues
	task.Creator = nil

	if err := h.taskRepo.Update(task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate task"})
		return
	}

	updated, _ := h.taskRepo.FindByID(id)
	c.JSON(http.StatusOK, gin.H{
		"message": "Task berhasil diupdate",
		"data":    updated,
	})
}

// UpdateStatus updates only the status of a task.
// PATCH /api/tasks/:id/status
func (h *TaskHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID task tidak valid"})
		return
	}

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status := models.TaskStatus(req.Status)
	if !status.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status tidak valid. Pilihan: todo, in_progress, done"})
		return
	}

	if err := h.taskRepo.UpdateStatus(id, status); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status berhasil diupdate"})
}

// Delete removes a task by ID.
// DELETE /api/tasks/:id
func (h *TaskHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID task tidak valid"})
		return
	}

	// Check task exists first
	if _, err := h.taskRepo.FindByID(id); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil task"})
		return
	}

	if err := h.taskRepo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task berhasil dihapus"})
}
