package repositories

import (
	"workms/internal/database"
	"workms/internal/models"

	"gorm.io/gorm"
)

type ProjectRepository struct {
	DB *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{DB: db}
}

// READ ALL
func (r *ProjectRepository) FindAll() ([]models.Project, error) {
	var proj []models.Project
	if err := database.DB.
		Find(&proj).Error; err != nil {
		return nil, err
	}
	return proj, nil
}

// CREATE
func (r *ProjectRepository) Create(proj *models.Project) error {
	return r.DB.Create(proj).Error
}

// READ ONE
func (r *ProjectRepository) FindByID(id uint64) (*models.Project, error) {
	var proj models.Project
	err := r.DB.First(&proj, id).Error
	return &proj, err
}

// UPDATE
func (r *ProjectRepository) Update(proj *models.Project) error {
	return r.DB.Save(proj).Error
}

// DELETE
func (r *ProjectRepository) Delete(id uint64) error {
	return r.DB.Delete(&models.Project{}, id).Error
}

// PAGINATION FUNCTION
func (r *ProjectRepository) FindPaginated(page, pageSize int) (tasks []models.Project, total int64, err error) {
	offset := (page - 1) * pageSize

	// total count
	if err = r.DB.Model(&models.Project{}).Count(&total).Error; err != nil {
		return
	}

	err = r.DB.
		Limit(pageSize).
		Offset(offset).
		Order("id asc").
		Find(&tasks).Error

	return
}

// GetProjectStatistics returns project counts
func (r *ProjectRepository) GetProjectStatistics() (total, active int64, err error) {
	err = r.DB.Model(&models.Project{}).Count(&total).Error
	if err != nil {
		return
	}

	err = r.DB.Model(&models.Project{}).Where("status = ?", "active").Count(&active).Error
	return
}

// FindAllWithRelations fetches all projects with their related Team and User
func (r *ProjectRepository) FindAllWithRelations() ([]models.Project, error) {
	var projects []models.Project
	err := r.DB.
		Preload("Team").
		Preload("User").
		Find(&projects).Error
	return projects, err
}

// GetTaskCountForProject returns the task count for a specific project
func (r *ProjectRepository) GetTaskCountForProject(projectID uint64) (total, completed int64, err error) {
	err = r.DB.Model(&models.Task{}).Where("project_id = ?", projectID).Count(&total).Error
	if err != nil {
		return
	}

	err = r.DB.Model(&models.Task{}).
		Where("project_id = ? AND status = ?", projectID, "completed").
		Count(&completed).Error
	return
}
