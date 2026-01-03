package models

import "time"

type Project struct {
	ID          uint64 `gorm:"primaryKey"`
	Name        string
	Description string
	Status      string // planned, active, completed

	StartDate time.Time
	EndDate   time.Time

	TeamID *uint64
	Team   *Team `gorm:"foreignKey:TeamID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	CreatedBy *uint64
	User      *User `gorm:"foreignKey:CreatedBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	CreatedAt time.Time
}
