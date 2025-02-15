package controllers

import (
	"log"
	"os"
	"strings"

	"notification_system/user_service/config"
	"notification_system/user_service/events"
	"notification_system/user_service/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

var kafkaProducer *events.KafkaProducer

func init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get the Kafka broker from environment variables
	kafkaBroker := os.Getenv("KAFKA_BROKER")
	if kafkaBroker == "" {
		log.Fatal("KAFKA_BROKER environment variable not set")
	}

	// Initialize Kafka producer
	kafkaProducer, err = events.NewKafkaProducer([]string{kafkaBroker}, "friendship_events")
	if err != nil {
		log.Fatal("Failed to create Kafka producer:", err)
	}
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

// Send Friend Request
func SendFriendRequest(c *fiber.Ctx) error {
	var input struct {
		Username string `json:"username"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var targetUser models.User
	result := config.DB.Where("username = ?", input.Username).First(&targetUser)
	if result.Error != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	// Check if the friendship request already exists
	var existingFriendship models.Friendship
	if err := config.DB.Where("(user_id1 = ? AND user_id2 = ?) OR (user_id1 = ? AND user_id2 = ?)", userID, targetUser.ID, targetUser.ID, userID).First(&existingFriendship).Error; err == nil {
		return c.Status(400).JSON(fiber.Map{"error": "Friend request already exists"})
	}

	// Create the friend request
	friendship := models.Friendship{UserID1: userID, UserID2: targetUser.ID, Status: "pending"}
	config.DB.Create(&friendship)

	// Publish Kafka event
	event := events.NotificationEvent{
		SenderID:   int(userID),
		ReceiverID: int(targetUser.ID),
		Message:    "You have received a new friend request",
		Topic:      "friend_request",
		Status:     "unread",
	}
	_ = kafkaProducer.SendMessage(event)

	return c.JSON(fiber.Map{"message": "Friend request sent"})
}

// View Friend Requests
func ViewFriendRequests(c *fiber.Ctx) error {
	userID, err := getUserIDFromToken(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var requests []models.Friendship
	config.DB.Where("user_id2 = ? AND status = ?", userID, "pending").Find(&requests)

	return c.JSON(requests)
}

// Accept Friend Request
func AcceptFriendRequest(c *fiber.Ctx) error {
	var input struct {
		Username string `json:"username"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var targetUser models.User
	result := config.DB.Where("username = ?", input.Username).First(&targetUser)
	if result.Error != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	// Check if the friend request is pending
	var friendship models.Friendship
	if err := config.DB.Where("user_id1 = ? AND user_id2 = ? AND status = ?", targetUser.ID, userID, "pending").First(&friendship).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "No pending friend request found"})
	}

	// Accept the friend request
	config.DB.Model(&friendship).Update("status", "accepted")

	// Publish Kafka event
	event := events.NotificationEvent{
		SenderID:   int(userID),
		ReceiverID: int(targetUser.ID),
		Message:    "Your friend request has been accepted",
		Topic:      "friend_request_accepted",
		Status:     "unread",
	}
	_ = kafkaProducer.SendMessage(event)

	return c.JSON(fiber.Map{"message": "Friend request accepted"})
}

// View Friends
func ViewFriends(c *fiber.Ctx) error {
	userID, err := getUserIDFromToken(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var friends []models.Friendship
	config.DB.Where("(user_id1 = ? OR user_id2 = ?) AND (status = ? OR status = ?)", userID, userID, "accepted", "unfollowed").Find(&friends)

	return c.JSON(friends)
}

// Unfriend
func Unfriend(c *fiber.Ctx) error {
	var input struct {
		Username string `json:"username"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var targetUser models.User
	result := config.DB.Where("username = ?", input.Username).First(&targetUser)
	if result.Error != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	config.DB.Where("(user_id1 = ? AND user_id2 = ?) OR (user_id1 = ? AND user_id2 = ?)", userID, targetUser.ID, targetUser.ID, userID).Delete(&models.Friendship{})

	return c.JSON(fiber.Map{"message": "User unfriended"})
}

// Unfollow
func Unfollow(c *fiber.Ctx) error {
	var input struct {
		Username string `json:"username"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var targetUser models.User
	result := config.DB.Where("username = ?", input.Username).First(&targetUser)
	if result.Error != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	config.DB.Model(&models.Friendship{}).Where("(user_id1 = ? AND user_id2 = ?) OR (user_id1 = ? AND user_id2 = ?) AND status = ?", userID, targetUser.ID, targetUser.ID, userID, "accepted").Update("status", "unfollowed")

	return c.JSON(fiber.Map{"message": "User unfollowed"})
}

// Follow Again
func FollowAgain(c *fiber.Ctx) error {
	var input struct {
		Username string `json:"username"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var targetUser models.User
	result := config.DB.Where("username = ?", input.Username).First(&targetUser)
	if result.Error != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	config.DB.Model(&models.Friendship{}).Where("(user_id1 = ? AND user_id2 = ?) OR (user_id1 = ? AND user_id2 = ?) AND status = ?", userID, targetUser.ID, targetUser.ID, userID, "unfollowed").Update("status", "accepted")

	return c.JSON(fiber.Map{"message": "User followed again"})
}
