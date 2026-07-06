package handlers

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/anisetusbambangmanalu19/moonlay-task-management-app/backend/internal/middleware"
	"github.com/anisetusbambangmanalu19/moonlay-task-management-app/backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	userRepo *repository.UserRepository
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(userRepo *repository.UserRepository) *AuthHandler {
	return &AuthHandler{userRepo: userRepo}
}

// LoginRequest defines the expected JSON body for POST /api/auth/login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=1"`
}

// Login handles user authentication and returns a JWT token.
// POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user by email
	user, err := h.userRepo.FindByEmail(req.Email)
	if err != nil {
		// Return generic message to avoid email enumeration
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah"})
		return
	}

	// Compare provided password with stored bcrypt hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah"})
		return
	}

	// Determine token expiry from env (default: 24 hours)
	expiryHours, _ := strconv.Atoi(os.Getenv("JWT_EXPIRY_HOURS"))
	if expiryHours <= 0 {
		expiryHours = 24
	}

	// Build JWT claims
	claims := &middleware.Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiryHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   strconv.FormatInt(user.ID, 10),
		},
	}

	// Sign token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenStr,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}
