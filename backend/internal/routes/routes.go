package routes

import (
	"github.com/anisetusbambangmanalu19/moonlay-task-management-app/backend/internal/handlers"
	"github.com/anisetusbambangmanalu19/moonlay-task-management-app/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

// Setup registers all routes for the application.
// Public routes: POST /api/auth/login
// Protected routes (JWT required): all others
func Setup(
	r *gin.Engine,
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	taskHandler *handlers.TaskHandler,
	chatbotHandler *handlers.ChatbotHandler,
) {
	api := r.Group("/api")

	// --- Public routes (no auth required) ---
	auth := api.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
	}

	// --- Protected routes (JWT required) ---
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware())
	{
		// Users
		protected.GET("/users", userHandler.GetAll)

		// Tasks
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
