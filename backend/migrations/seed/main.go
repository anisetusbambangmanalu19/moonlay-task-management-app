package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// User model - minimal struct just for seeding
type User struct {
	ID           int64     `gorm:"primaryKey;column:id"`
	Name         string    `gorm:"column:name"`
	Email        string    `gorm:"column:email"`
	PasswordHash string    `gorm:"column:password_hash"`
	CreatedAt    time.Time `gorm:"column:created_at"`
}

func (User) TableName() string { return "users" }

// hashPassword generates a bcrypt hash from plaintext password
func hashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("❌ Gagal hash password: %v", err)
	}
	return string(hash)
}

func main() {
	// Load .env from backend root (run this from backend/ directory)
	if err := godotenv.Load(".env"); err != nil {
		log.Println("⚠️  Tidak ada file .env, membaca dari environment system")
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "postgres"),
		os.Getenv("DB_PASSWORD"),
		getEnv("DB_NAME", "moonlay_task_db"),
		getEnv("DB_SSLMODE", "disable"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Gagal koneksi database: %v", err)
	}

	log.Println("✅ Terhubung ke database")

	// Seed users with bcrypt-hashed passwords
	// Passwords are hashed at runtime — never stored as plaintext
	seedUsers := []struct {
		Name     string
		Email    string
		Password string
	}{
		{"Admin", "admin@moonlay.com", "admin123"},
		{"Budi Santoso", "budi@moonlay.com", "budi123"},
		{"Siti Rahayu", "siti@moonlay.com", "siti123"},
		{"Rangga Pratama", "rangga@moonlay.com", "rangga123"},
	}

	for _, su := range seedUsers {
		user := User{
			Name:         su.Name,
			Email:        su.Email,
			PasswordHash: hashPassword(su.Password),
		}

		// Use FirstOrCreate to avoid duplicate insertion on re-run
		result := db.Where(User{Email: su.Email}).FirstOrCreate(&user)
		if result.Error != nil {
			log.Printf("❌ Gagal seed user %s: %v", su.Email, result.Error)
		} else if result.RowsAffected > 0 {
			log.Printf("✅ User dibuat: %s (%s)", su.Name, su.Email)
		} else {
			log.Printf("⚠️  User sudah ada (skip): %s (%s)", su.Name, su.Email)
		}
	}

	log.Println("\n🎉 Seed selesai!")
	log.Println("Kredensial login yang tersedia:")
	log.Println("  admin@moonlay.com    / admin123")
	log.Println("  budi@moonlay.com     / budi123")
	log.Println("  siti@moonlay.com     / siti123")
	log.Println("  rangga@moonlay.com   / rangga123")
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
