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

// Model User - struct minimal khusus untuk keperluan seeding
type User struct {
	ID           int64     `gorm:"primaryKey;column:id"`
	Name         string    `gorm:"column:name"`
	Email        string    `gorm:"column:email"`
	PasswordHash string    `gorm:"column:password_hash"`
	CreatedAt    time.Time `gorm:"column:created_at"`
}

func (User) TableName() string { return "users" }

// hashPassword membuat hash bcrypt dari password plaintext
func hashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("❌ Gagal hash password: %v", err)
	}
	return string(hash)
}

func main() {
	// Muat .env dari root backend (jalankan dari direktori backend/)
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

	// Seed user dengan password yang sudah di-hash bcrypt
	// Password di-hash saat runtime — tidak pernah disimpan plaintext
	seedUsers := []struct {
		Name     string
		Email    string
		Password string
	}{
		{"Admin", "admin@moonlay.com", "admin123"},
		{"Anisetus", "anisetus@moonlay.com", "anisetus123"},
		{"Bambang", "bambang@moonlay.com", "bambang123"},
		{"Manalu", "manalu@moonlay.com", "manalu123"},
	}

	for _, su := range seedUsers {
		user := User{
			Name:         su.Name,
			Email:        su.Email,
			PasswordHash: hashPassword(su.Password),
		}

		// Gunakan FirstOrCreate agar tidak membuat data ganda saat dijalankan ulang
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
	log.Println("  anisetus@moonlay.com / anisetus123")
	log.Println("  bambang@moonlay.com  / bambang123")
	log.Println("  manalu@moonlay.com   / manalu123")
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
