package handlers

import (
	"net/http"

	"github.com/anisetusbambangmanalu19/moonlay-task-management-app/backend/internal/repository"
	"github.com/gin-gonic/gin"
)

// UserHandler menangani permintaan HTTP terkait user
type UserHandler struct {
	userRepo *repository.UserRepository
}

// NewUserHandler membuat UserHandler baru
func NewUserHandler(userRepo *repository.UserRepository) *UserHandler {
	return &UserHandler{userRepo: userRepo}
}

// GetAll mengembalikan daftar semua user (untuk dropdown assignee di frontend).
// GET /api/users
func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.userRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users})
}
