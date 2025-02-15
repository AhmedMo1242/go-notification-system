package models

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"type:varchar(255);uniqueIndex"`
	Email     string `gorm:"type:varchar(255);uniqueIndex"`
	Password  string `gorm:"type:varchar(255)"`
	CreatedAt time.Time
	LastLogin time.Time
}
