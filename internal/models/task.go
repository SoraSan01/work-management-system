package models

import "time"

type Task struct {
	ID          uint64 `gorm:"primaryKey"`
	Title       string
	Description string

	Status   string // todo, in_progress, done
	Priority string // low, medium, high

	ProjectID *uint64
	Project   *Project `gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	AssignedTo *uint64
	User       *User `gorm:"foreignKey:AssignedTo;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	DueDate   time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
