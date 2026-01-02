package models

type Role struct {
	ID    uint64 `gorm:"primaryKey;autoIncrement"`
	Name  string `gorm:"size:50;unique;not null"`
	Users []User `gorm:"foreignKey:RoleID"`
}
