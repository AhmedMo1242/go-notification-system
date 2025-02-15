package models

import "time"

type Notification struct {
	ID         uint      `gorm:"primaryKey"`
	SenderID   int       `json:"sender_id"`
	ReceiverID int       `json:"receiver_id"`
	Message    string    `json:"message"`
	Topic      string    `json:"topic"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}
