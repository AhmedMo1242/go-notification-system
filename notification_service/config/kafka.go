package config

import (
	"log"
	"os"

	"github.com/IBM/sarama"
	"github.com/joho/godotenv"
)

var KafkaProducer sarama.SyncProducer
var KafkaConsumer sarama.Consumer

func InitKafka() {
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

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{kafkaBroker}, config)
	if err != nil {
		log.Fatal("Failed to start Kafka producer:", err)
	}

	consumer, err := sarama.NewConsumer([]string{kafkaBroker}, config)
	if err != nil {
		log.Fatal("Failed to start Kafka consumer:", err)
	}

	KafkaProducer = producer
	KafkaConsumer = consumer
	log.Println("Kafka producer and consumer started")
}
