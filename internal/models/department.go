package models

import "time"

type Department struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"size:100;not null;unique"` // Department name must be unique
	IsActive  bool   `gorm:"not null"`                 // Optional: mark department as active/inactive
	Users     []User `gorm:"foreignKey:DepartmentID"`  // One-to-many relationship with User
	CreatedAt time.Time
	UpdatedAt time.Time
}
