package models

import "time"

type Task struct {
	ID          uint64 `gorm:"primaryKey"`
	Title       string
	Description string

	Status   string
	Priority string

	ProjectID *uint64
	Project   *Project

	AssignedTo *uint64
	User       *User `gorm:"foreignKey:AssignedTo;references:ID"`

	DueDate   time.Time
	CreatedAt time.Time
	UpdatedAt time.Time

	Documents []Document `gorm:"foreignKey:TaskID"`
}
