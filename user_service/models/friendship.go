package models

import "time"

type Friendship struct {
	UserID1   uint `gorm:"primaryKey"`
	UserID2   uint `gorm:"primaryKey"`
	Status    string
	CreatedAt time.Time
}
