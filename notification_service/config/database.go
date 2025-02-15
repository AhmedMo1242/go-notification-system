package config

import (
	"log"
	"notification_system/notification_service/models"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get the DSN from environment variables
	dsn := os.Getenv("DB_DSN_NOTIFICATION")
	if dsn == "" {
		log.Fatal("DB_DSN_NOTIFICATION environment variable not set")
	}

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB.AutoMigrate(&models.Notification{})
}
