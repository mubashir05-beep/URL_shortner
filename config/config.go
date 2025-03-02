package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mubashir05-beep/url_shortner/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Global DB variable
var DB *gorm.DB

func ConnectDB() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get DATABASE_URL
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not found in .env")
	}

	// Connect to PostgreSQL with GORM
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run AutoMigrations
	err = db.AutoMigrate(&models.User{}, &models.URL{}, &models.Analytics{})
	if err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Set DB instance
	DB = db

	fmt.Println("Connected to Database!")
}
