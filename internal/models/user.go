package models

import "time"

type User struct {
	ID           uint64     `gorm:"primaryKey;autoIncrement"`
	FirstName    string     `gorm:"size:50;not null"`
	LastName     string     `gorm:"size:50;not null"`
	Email        string     `gorm:"size:100;not null;unique"`
	Password     string     `gorm:"not null"`
	RoleID       uint64     `gorm:"column:role_id"`
	Role         Role       `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	IsActive     bool       `gorm:"not null"`
	DepartmentID uint64     `gorm:"column:department_id"`
	Department   Department `gorm:"foreignKey:DepartmentID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
	