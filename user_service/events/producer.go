package events

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
)

// NotificationEvent represents the notification structure for Kafka
type NotificationEvent struct {
	SenderID   int    `json:"sender_id"`
	ReceiverID int    `json:"receiver_id"`
	Message    string `json:"message"`
	Topic      string `json:"topic"`
	Status     string `json:"status"`
}

// KafkaProducer holds the producer instance
type KafkaProducer struct {
	producer sarama.SyncProducer
	topic    string
}

// NewKafkaProducer initializes a new Kafka producer
func NewKafkaProducer(brokers []string, topic string) (*KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{producer: producer, topic: topic}, nil
}

// SendMessage sends a structured notification event to Kafka
func (p *KafkaProducer) SendMessage(event NotificationEvent) error {
	messageBytes, err := json.Marshal(event)
	if err != nil {
		log.Println("Failed to serialize message:", err)
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(messageBytes),
	}

	_, _, err = p.producer.SendMessage(msg)
	if err != nil {
		log.Println("Failed to send Kafka message:", err)
		return err
	}

	return nil
}

// Close shuts down the producer
func (p *KafkaProducer) Close() error {
	return p.producer.Close()
}
