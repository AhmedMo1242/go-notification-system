package main

import (
	"log"
	"os"

	"notification_system/notification_service/config"
	"notification_system/notification_service/events"
	"notification_system/notification_service/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get the notification service port from environment variables
	port := os.Getenv("NOTIFICATION_SERVICE_PORT")
	if port == "" {
		log.Fatal("NOTIFICATION_SERVICE_PORT environment variable not set")
	}

	// Get the Kafka broker from environment variables
	kafkaBroker := os.Getenv("KAFKA_BROKER")
	if kafkaBroker == "" {
		log.Fatal("KAFKA_BROKER environment variable not set")
	}

	// Initialize database
	config.ConnectDatabase()

	// Start Kafka consumer in a separate goroutine
	go events.StartKafkaConsumer([]string{kafkaBroker}, "friendship_events")

	// Initialize Fiber app
	app := fiber.New()

	// Register API routes
	routes.RegisterRoutes(app)

	log.Println("notification service running on port", port)

	// Start Fiber server
	log.Fatal(app.Listen(":" + port))
}
