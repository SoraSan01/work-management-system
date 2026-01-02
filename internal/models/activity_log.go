package models

import "time"

type ActivityLog struct {
	ID        uint64 `gorm:"primaryKey"`
	UserID    uint64
	Action    string // created_task, updated_task, completed_task
	Entity    string // task, project, document
	EntityID  uint64
	CreatedAt time.Time
}
