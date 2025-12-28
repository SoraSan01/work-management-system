package models

import "time"

type Team struct {
	ID          uint64 `gorm:"primaryKey"`
	Name        string
	Description string
	CreatedBy   uint // User ID
	CreatedAt   time.Time
}
