package handlers

import (
	"net/http"

	"github.com/anisetusbambangmanalu19/moonlay-task-management-app/backend/internal/repository"
	"github.com/gin-gonic/gin"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userRepo *repository.UserRepository
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userRepo *repository.UserRepository) *UserHandler {
	return &UserHandler{userRepo: userRepo}
}

// GetAll returns a list of all users (for assignee dropdown in frontend).
// GET /api/users
func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.userRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users})
}
