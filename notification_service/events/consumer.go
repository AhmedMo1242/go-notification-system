package events

import (
	"encoding/json"
	"log"
	"notification_system/notification_service/config"
	"notification_system/notification_service/models"

	"github.com/IBM/sarama"
)

// NotificationEvent matches the structure used in Kafka messages
type NotificationEvent struct {
	SenderID   int    `json:"sender_id"`
	ReceiverID int    `json:"receiver_id"`
	Message    string `json:"message"`
	Topic      string `json:"topic"`
	Status     string `json:"status"`
}

// StartKafkaConsumer listens for Kafka messages and saves them as notifications
func StartKafkaConsumer(brokers []string, topic string) {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokers, kafkaConfig)
	if err != nil {
		log.Fatal("Failed to start consumer:", err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatal("Failed to consume partition:", err)
	}
	defer partitionConsumer.Close()

	for msg := range partitionConsumer.Messages() {
		var event NotificationEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Println("Failed to parse Kafka message:", err)
			continue
		}

		// Save the notification in the database
		notification := models.Notification{
			SenderID:   event.SenderID,
			ReceiverID: event.ReceiverID,
			Message:    event.Message,
			Topic:      event.Topic,
			Status:     event.Status,
		}

		config.DB.Create(&notification)
		log.Println("Stored notification:", notification)
	}
}
