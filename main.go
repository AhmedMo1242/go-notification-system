package main

import (
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get the microservice endpoints from environment variables
	userService := os.Getenv("USER_SERVICE_URL")
	if userService == "" {
		log.Fatal("USER_SERVICE_URL environment variable not set")
	}

	notificationService := os.Getenv("NOTIFICATION_SERVICE_URL")
	if notificationService == "" {
		log.Fatal("NOTIFICATION_SERVICE_URL environment variable not set")
	}

	// Get the API Gateway port from environment variables
	gatewayPort := os.Getenv("API_GATEWAY_PORT")
	if gatewayPort == "" {
		log.Fatal("API_GATEWAY_PORT environment variable not set")
	}

	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		path := c.Path()

		if strings.HasPrefix(path, "/user") {
			// Forward to user_service
			return proxy.Do(c, userService+path+"?"+c.Request().URI().QueryArgs().String())
		} else if strings.HasPrefix(path, "/notification") {
			// Forward to notification_service
			return proxy.Do(c, notificationService+path+"?"+c.Request().URI().QueryArgs().String())
		}

		// If no match, return 404
		return c.Status(404).JSON(fiber.Map{"error": "Service not found"})
	})

	log.Println("API Gateway running on port", gatewayPort)
	log.Fatal(app.Listen(":" + gatewayPort))
}
