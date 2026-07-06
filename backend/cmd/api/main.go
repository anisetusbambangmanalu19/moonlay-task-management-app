package main

import (
	"log"
	"os"

	"github.com/anisetusbambangmanalu19/moonlay-task-management-app/backend/internal/config"
	"github.com/anisetusbambangmanalu19/moonlay-task-management-app/backend/internal/handlers"
	"github.com/anisetusbambangmanalu19/moonlay-task-management-app/backend/internal/repository"
	"github.com/anisetusbambangmanalu19/moonlay-task-management-app/backend/internal/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load environment variables from .env file
	config.LoadEnv()

	// Connect to PostgreSQL database
	config.ConnectDB()
	db := config.GetDB()

	// Initialize repositories (query layer)
	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)

	// Initialize handlers (HTTP layer)
	authHandler := handlers.NewAuthHandler(userRepo)
	userHandler := handlers.NewUserHandler(userRepo)
	taskHandler := handlers.NewTaskHandler(taskRepo, userRepo)
	chatbotHandler := handlers.NewChatbotHandler(taskRepo)

	// Setup Gin router
	r := gin.Default()

	// CORS middleware — allow frontend (Next.js on :3000) to communicate with backend
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Register all application routes
	routes.Setup(r, authHandler, userHandler, taskHandler, chatbotHandler)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server berjalan di http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
