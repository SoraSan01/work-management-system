package models

import "time"

type Notification struct {
	ID        uint64 `gorm:"primaryKey"`
	UserID    uint64
	Title     string
	Message   string
	IsRead    bool
	CreatedAt time.Time
}
