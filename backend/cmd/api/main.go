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
	// Muat variabel lingkungan dari file .env
	config.LoadEnv()

	// Hubungkan ke database PostgreSQL
	config.ConnectDB()
	db := config.GetDB()

	// Inisialisasi repository (lapisan query)
	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)

	// Inisialisasi handler (lapisan HTTP)
	authHandler := handlers.NewAuthHandler(userRepo)
	userHandler := handlers.NewUserHandler(userRepo)
	taskHandler := handlers.NewTaskHandler(taskRepo, userRepo)
	chatbotHandler := handlers.NewChatbotHandler(taskRepo)

	// Siapkan router Gin
	r := gin.Default()

	// Middleware CORS untuk mengizinkan frontend (Next.js di :3000) berkomunikasi dengan backend
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Daftarkan semua rute aplikasi
	routes.Setup(r, authHandler, userHandler, taskHandler, chatbotHandler)

	// Jalankan server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server berjalan di http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
