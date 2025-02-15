package controllers

import (
	"log"
	"os"
	"strings"
	"time"

	"notification_system/user_service/config"
	"notification_system/user_service/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

var jwtSecret []byte
var blacklistedTokens = make(map[string]time.Time)

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

// Signup (Register a new user)
func Signup(c *fiber.Ctx) error {
	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Check if username or email already exists
	var existingUser models.User
	if err := config.DB.Where("username = ? OR email = ?", input.Username, input.Email).First(&existingUser).Error; err == nil {
		return c.Status(400).JSON(fiber.Map{"error": "Username or email already taken"})
	}

	// Create new user
	user := models.User{
		Username:  input.Username,
		Email:     input.Email,
		Password:  input.Password,
		CreatedAt: time.Now(),
		LastLogin: time.Now(),
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create user"})
	}

	return c.JSON(fiber.Map{"message": "User created successfully"})
}

// Login with JWT handling
func Login(c *fiber.Ctx) error {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	var user models.User
	result := config.DB.Where("username = ? AND password = ?", input.Username, input.Password).First(&user)
	if result.Error != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid username or password"})
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.ID,
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	user.LastLogin = time.Now()
	config.DB.Save(&user)

	return c.JSON(fiber.Map{"message": "Login successful", "token": tokenString})
}

// Logout with JWT handling and token blacklisting
func Logout(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Blacklist the token
	blacklistedTokens[tokenString] = time.Now().Add(72 * time.Hour)

	return c.JSON(fiber.Map{"message": "Logged out successfully"})
}

// Middleware to check if token is blacklisted
func CheckBlacklist(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Check if token is blacklisted
	if exp, found := blacklistedTokens[tokenString]; found && exp.After(time.Now()) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token is blacklisted"})
	}

	return c.Next()
}
