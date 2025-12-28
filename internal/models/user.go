package models

import "time"

type User struct {
	ID        uint64 `gorm:"primaryKey"`
	FirstName string
	LastName  string
	Email     string `gorm:"unique"`
	Password  string
	Role      string // admin, manager, employee
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
