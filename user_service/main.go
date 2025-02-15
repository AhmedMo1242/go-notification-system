package main

import (
	"log"
	"notification_system/user_service/config"
	"notification_system/user_service/routes"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Initialize Fiber
	app := fiber.New()

	// Connect to Database
	config.ConnectDatabase()

	// Register Routes
	routes.AuthRoutes(app)
	routes.FriendshipRoutes(app)

	// Start Server
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("USER_SERVICE_PORT")
	if port == "" {
		log.Fatal("USER_SERVICE_PORT environment variable not set")
	}

	log.Println("user server running on port", port)

	log.Fatal(app.Listen(":" + port))
}
