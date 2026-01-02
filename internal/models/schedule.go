package models

import "time"

type Schedule struct {
	ID        uint64 `gorm:"primaryKey"`
	UserID    uint64
	TaskID    uint64
	StartTime time.Time
	EndTime   time.Time
	Type      string // meeting, task, reminder
}
