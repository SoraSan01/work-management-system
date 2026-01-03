package models

import "time"

type Team struct {
	ID          uint64 `gorm:"primaryKey"`
	Name        string
	Description string

	CreatedBy *uint64
	User      *User `gorm:"foreignKey:CreatedBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	Members []TeamMember `gorm:"foreignKey:TeamID"`

	CreatedAt time.Time
}
