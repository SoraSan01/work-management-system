package repositories

import (
	"workms/internal/models"

	"gorm.io/gorm"
)

type DocumentRepository struct {
	DB *gorm.DB
}

func NewDocumentRepository(db *gorm.DB) *DocumentRepository {
	return &DocumentRepository{DB: db}
}

// Create a new document
func (r *DocumentRepository) Create(doc *models.Document) error {
	return r.DB.Create(doc).Error
}

// FindAll retrieves all documents with relations
func (r *DocumentRepository) FindAll() ([]models.Document, error) {
	var documents []models.Document
	err := r.DB.
		Preload("User").
		Preload("Project").
		Preload("Project.Team").
		Preload("Task").
		Where("is_deleted = ?", false).
		Order("created_at DESC").
		Find(&documents).Error
	return documents, err
}

// FindByID retrieves a single document by ID
func (r *DocumentRepository) FindByID(id uint64) (*models.Document, error) {
	var document models.Document
	err := r.DB.
		Preload("User").
		Preload("Project").
		Preload("Task").
		Where("id = ? AND is_deleted = ?", id, false).
		First(&document).Error
	return &document, err
}

// Update updates a document
func (r *DocumentRepository) Update(doc *models.Document) error {
	return r.DB.Save(doc).Error
}

// Delete soft deletes a document
func (r *DocumentRepository) Delete(id uint64) error {
	return r.DB.Model(&models.Document{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"deleted_at": gorm.Expr("NOW()"),
		}).Error
}

// HardDelete permanently deletes a document
func (r *DocumentRepository) HardDelete(id uint64) error {
	return r.DB.Unscoped().Delete(&models.Document{}, id).Error
}

// GetDocumentStatistics returns document statistics
func (r *DocumentRepository) GetDocumentStatistics() (total int64, totalSize int64, recentUploads int64, fileTypes int64, err error) {
	// Total documents
	err = r.DB.Model(&models.Document{}).
		Where("is_deleted = ?", false).
		Count(&total).Error
	if err != nil {
		return
	}

	// Total size
	var result struct {
		TotalSize int64
	}
	err = r.DB.Model(&models.Document{}).
		Where("is_deleted = ?", false).
		Select("COALESCE(SUM(file_size), 0) as total_size").
		Scan(&result).Error
	if err != nil {
		return
	}
	totalSize = result.TotalSize

	// Recent uploads (last 7 days)
	err = r.DB.Model(&models.Document{}).
		Where("is_deleted = ? AND created_at >= NOW() - INTERVAL '7 days'", false).
		Count(&recentUploads).Error
	if err != nil {
		return
	}

	// Distinct file types
	err = r.DB.Model(&models.Document{}).
		Where("is_deleted = ?", false).
		Distinct("file_type").
		Count(&fileTypes).Error

	return
}

// FindByProjectID retrieves documents for a specific project
func (r *DocumentRepository) FindByProjectID(projectID uint64) ([]models.Document, error) {
	var documents []models.Document
	err := r.DB.
		Preload("User").
		Preload("Project").
		Preload("Task").
		Where("project_id = ? AND is_deleted = ?", projectID, false).
		Order("created_at DESC").
		Find(&documents).Error
	return documents, err
}

// FindByTaskID retrieves documents for a specific task
func (r *DocumentRepository) FindByTaskID(taskID uint64) ([]models.Document, error) {
	var documents []models.Document
	err := r.DB.
		Preload("User").
		Preload("Project").
		Preload("Task").
		Where("task_id = ? AND is_deleted = ?", taskID, false).
		Order("created_at DESC").
		Find(&documents).Error
	return documents, err
}

// FindByUserID retrieves documents uploaded by a specific user
func (r *DocumentRepository) FindByUserID(userID uint64) ([]models.Document, error) {
	var documents []models.Document
	err := r.DB.
		Preload("User").
		Preload("Project").
		Preload("Task").
		Where("uploaded_by = ? AND is_deleted = ?", userID, false).
		Order("created_at DESC").
		Find(&documents).Error
	return documents, err
}

// Search documents by name
func (r *DocumentRepository) Search(query string) ([]models.Document, error) {
	var documents []models.Document
	err := r.DB.
		Preload("User").
		Preload("Project").
		Preload("Task").
		Where("(name ILIKE ? OR original_name ILIKE ?) AND is_deleted = ?", "%"+query+"%", "%"+query+"%", false).
		Order("created_at DESC").
		Find(&documents).Error
	return documents, err
}
