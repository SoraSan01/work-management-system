package models

type TeamMember struct {
	ID     uint64 `gorm:"primaryKey"`
	TeamID uint64
	UserID uint64
	Role   string // leader, member
}
