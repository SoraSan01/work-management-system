package models

import "time"

type Document struct {
	ID         uint64 `gorm:"primaryKey"`
	Name       string
	FilePath   string
	FileType   string
	UploadedBy uint64
	ProjectID  uint64
	TaskID     uint64
	CreatedAt  time.Time
}
