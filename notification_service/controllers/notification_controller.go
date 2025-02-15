package controllers

import (
	"log"
	"os"
	"strings"

	"notification_system/notification_service/config"
	"notification_system/notification_service/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

var jwtSecret []byte

func init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get the JWT secret from environment variables
	jwtSecretEnv := os.Getenv("JWT_SECRET")
	if jwtSecretEnv == "" {
		log.Fatal("JWT_SECRET environment variable not set")
	}
	jwtSecret = []byte(jwtSecretEnv)
}

// Helper function to retrieve user ID from JWT token
func getUserIDFromToken(c *fiber.Ctx) (uint, error) {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return 0, fiber.ErrUnauthorized
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.ErrUnauthorized
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return 0, fiber.ErrUnauthorized
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, fiber.ErrUnauthorized
	}

	userID, ok := claims["userID"].(float64)
	if !ok {
		return 0, fiber.ErrUnauthorized
	}

	return uint(userID), nil
}

// Get all notifications for a user
func GetNotifications(c *fiber.Ctx) error {
	userID, err := getUserIDFromToken(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var notifications []models.Notification
	if err := config.DB.Where("receiver_id = ?", userID).Find(&notifications).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch notifications"})
	}

	return c.JSON(notifications)
}

// Mark a notification as read
func MarkAsRead(c *fiber.Ctx) error {
	notificationID := c.Params("notification_id")
	if notificationID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Notification ID is required"})
	}

	var notification models.Notification
	if err := config.DB.First(&notification, notificationID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Notification not found"})
	}

	notification.Status = "read"
	if err := config.DB.Save(&notification).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update notification"})
	}

	return c.JSON(fiber.Map{"message": "Notification marked as read"})
}
