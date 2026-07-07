package routes

import (
	"github.com/anisetusbambangmanalu19/moonlay-task-management-app/backend/internal/handlers"
	"github.com/anisetusbambangmanalu19/moonlay-task-management-app/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

// Setup mendaftarkan semua rute untuk aplikasi.
// Rute publik: POST /api/auth/login
// Rute terlindungi (wajib JWT): sisanya
func Setup(
	r *gin.Engine,
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	taskHandler *handlers.TaskHandler,
	chatbotHandler *handlers.ChatbotHandler,
) {
	api := r.Group("/api")

	// --- Rute publik (tanpa autentikasi) ---
	auth := api.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
	}

	// --- Rute terlindungi (wajib JWT) ---
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware())
	{
		// User
		protected.GET("/users", userHandler.GetAll)

		// Task
		tasks := protected.Group("/tasks")
		{
			tasks.GET("", taskHandler.GetAll)
			tasks.GET("/:id", taskHandler.GetByID)
			tasks.POST("", taskHandler.Create)
			tasks.PUT("/:id", taskHandler.Update)
			tasks.PATCH("/:id/status", taskHandler.UpdateStatus)
			tasks.DELETE("/:id", taskHandler.Delete)
		}

		// Chatbot
		protected.POST("/chatbot", chatbotHandler.Chat)
	}
}
