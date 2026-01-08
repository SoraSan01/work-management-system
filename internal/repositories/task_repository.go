package repositories

import (
	"workms/internal/database"
	"workms/internal/models"

	"gorm.io/gorm"
)

type TaskRepository struct {
	DB *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{DB: db}
}

// READ ALL
func (r *TaskRepository) FindAll() ([]models.Task, error) {
	var task []models.Task
	if err := database.DB.
		Find(&task).Error; err != nil {
		return nil, err
	}
	return task, nil
}

// CREATE
func (r *TaskRepository) Create(proj *models.Task) error {
	return r.DB.Create(proj).Error
}

// READ ONE
func (r *TaskRepository) FindByID(id uint64) (*models.Task, error) {
	var task models.Task
	err := r.DB.First(&task, id).Error
	return &task, err
}

// UPDATE
func (r *TaskRepository) Update(proj *models.Task) error {
	return r.DB.Save(proj).Error
}

// DELETE
func (r *TaskRepository) Delete(id uint64) error {
	return r.DB.Delete(&models.Task{}, id).Error
}

// PAGINATION FUNCTION
func (r *TaskRepository) FindPaginated(page, pageSize int) (tasks []models.Task, total int64, err error) {
	offset := (page - 1) * pageSize

	// total count
	if err = r.DB.Model(&models.Task{}).Count(&total).Error; err != nil {
		return
	}

	err = r.DB.
		Preload("Project").
		Preload("User").
		Limit(pageSize).
		Offset(offset).
		Order("id asc").
		Find(&tasks).Error

	return
}

// FindAllWithRelations fetches all tasks with their related Project and User
func (r *TaskRepository) FindAllWithRelations() ([]models.Task, error) {
	var tasks []models.Task
	err := r.DB.
		Preload("Project").
		Preload("Project.Team").
		Preload("User").
		Find(&tasks).Error
	return tasks, err
}

// GetTaskStatistics returns task count by status
func (r *TaskRepository) GetTaskStatistics() (total, inProgress, completed int64, err error) {
	err = r.DB.Model(&models.Task{}).Count(&total).Error
	if err != nil {
		return
	}

	err = r.DB.Model(&models.Task{}).Where("status = ?", "in_progress").Count(&inProgress).Error
	if err != nil {
		return
	}

	err = r.DB.Model(&models.Task{}).Where("status = ?", "done").Count(&completed).Error
	return
}
