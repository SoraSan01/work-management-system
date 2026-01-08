package models

import (
	"time"
)

type Document struct {
	ID           uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string `gorm:"type:varchar(255);not null" json:"name"`
	OriginalName string `gorm:"type:varchar(255);not null" json:"original_name"`
	FilePath     string `gorm:"type:varchar(500);not null" json:"file_path"`
	FileType     string `gorm:"type:varchar(100)" json:"file_type"`
	FileSize     int64  `gorm:"not null" json:"file_size"`
	MimeType     string `gorm:"type:varchar(100)" json:"mime_type"`

	UploadedBy *uint64 `gorm:"index" json:"uploaded_by"`
	ProjectID  *uint64 `gorm:"index" json:"project_id"`
	TaskID     *uint64 `gorm:"index" json:"task_id"`

	User    *User    `gorm:"foreignKey:UploadedBy;references:ID" json:"user,omitempty"`
	Project *Project `gorm:"foreignKey:ProjectID;references:ID" json:"project,omitempty"`
	Task    *Task    `gorm:"foreignKey:TaskID;references:ID" json:"task,omitempty"`

	Description string `gorm:"type:text" json:"description"`
	Version     int    `gorm:"default:1" json:"version"`
	IsDeleted   bool   `gorm:"default:false;index" json:"is_deleted"`

	// Approval Workflow
	IsApproved bool       `gorm:"default:false" json:"is_approved"`
	ApprovedBy *uint64    `gorm:"index" json:"approved_by"`
	ApprovedAt *time.Time `json:"approved_at,omitempty"`

	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}
