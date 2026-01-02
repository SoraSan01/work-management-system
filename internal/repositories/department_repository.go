package repositories

import (
	"workms/internal/models"

	"gorm.io/gorm"
)

type DepartmentRepository struct {
	DB *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) *DepartmentRepository {
	return &DepartmentRepository{DB: db}
}

// CREATE
func (r *DepartmentRepository) Create(department *models.Department) error {
	return r.DB.Create(department).Error
}

// READ ALL
func (r *DepartmentRepository) FindAll() ([]models.Department, error) {
	var department []models.Department
	err := r.DB.Find(&department).Error
	return department, err
}

// READ ONE
func (r *DepartmentRepository) FindByID(id uint64) (*models.Department, error) {
	var department models.Department
	err := r.DB.First(&department, id).Error
	return &department, err
}

// UPDATE
func (r *DepartmentRepository) Update(department *models.Department) error {
	return r.DB.Save(department).Error
}

// DELETE
func (r *DepartmentRepository) Delete(id uint64) error {
	return r.DB.Delete(&models.Department{}, id).Error
}

// PAGINATION FUNCTION
func (r *DepartmentRepository) FindPaginated(page, pageSize int) (users []models.Department, total int64, err error) {
	offset := (page - 1) * pageSize

	// total count
	if err = r.DB.Model(&models.Department{}).Count(&total).Error; err != nil {
		return
	}

	err = r.DB.
		Limit(pageSize).
		Offset(offset).
		Order("id asc").
		Find(&users).Error

	return
}
