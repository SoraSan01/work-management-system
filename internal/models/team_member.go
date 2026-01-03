package models

type TeamMember struct {
	ID uint64 `gorm:"primaryKey"`

	TeamID uint64
	Team   Team `gorm:"foreignKey:TeamID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	UserID uint64
	User   User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	Role string // leader, member
}
