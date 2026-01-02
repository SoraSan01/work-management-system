package models

import "time"

type Task struct {
	ID          uint64 `gorm:"primaryKey"`
	Title       string
	Description string
	Status      string // todo, in_progress, done
	Priority    string // low, medium, high
	ProjectID   uint64
	AssignedTo  uint64 // User ID
	DueDate     time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
