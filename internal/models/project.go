package models

import "time"

type Project struct {
	ID          uint64 `gorm:"primaryKey"`
	Name        string
	Description string
	Status      string // planned, active, completed
	StartDate   time.Time
	EndDate     time.Time
	TeamID      uint64
	CreatedBy   uint64
	CreatedAt   time.Time
}
