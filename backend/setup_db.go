package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: Error loading .env file")
	}

	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5434")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgrespassword")
	dbname := getEnv("DB_NAME", "moonlay_task_db")

	// 1. Connect to postgres database (default) to create our database
	connStrDefault := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
		host, port, user, password)
	
	dbDefault, err := sql.Open("postgres", connStrDefault)
	if err != nil {
		log.Fatalf("❌ Failed to connect to default postgres db: %v", err)
	}
	defer dbDefault.Close()

	// Check if database exists
	var exists bool
	err = dbDefault.QueryRow("SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = $1)", dbname).Scan(&exists)
	if err != nil {
		log.Fatalf("❌ Failed to check if database exists: %v", err)
	}

	if !exists {
		log.Printf("Creating database %s...", dbname)
		_, err = dbDefault.Exec(fmt.Sprintf("CREATE DATABASE %s", dbname))
		if err != nil {
			log.Fatalf("❌ Failed to create database: %v", err)
		}
		log.Println("✅ Database created successfully!")
	} else {
		log.Printf("✅ Database %s already exists.", dbname)
	}

	// 2. Connect to our new database to run migrations
	connStrNew := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	
	dbNew, err := sql.Open("postgres", connStrNew)
	if err != nil {
		log.Fatalf("❌ Failed to connect to new database: %v", err)
	}
	defer dbNew.Close()

	// 3. Run migration
	migrationPath := filepath.Join("migrations", "001_init.sql")
	migrationData, err := os.ReadFile(migrationPath)
	if err != nil {
		log.Fatalf("❌ Failed to read migration file: %v", err)
	}

	log.Println("Running migration 001_init.sql...")
	_, err = dbNew.Exec(string(migrationData))
	if err != nil {
		log.Fatalf("❌ Failed to execute migration: %v", err)
	}
	log.Println("✅ Migration executed successfully!")
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
