package repositories

import (
	"errors"
	"workms/internal/database"
	"workms/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// CREATE
func (r *UserRepository) Create(user *models.User) error {
	return r.DB.Create(user).Error
}

// GET ACTIVE DEPARTMENTS
func (r *UserRepository) GetActiveDepartments() ([]models.Department, error) {
	var depts []models.Department
	if err := r.DB.Where("is_active = ?", true).Find(&depts).Error; err != nil {
		return nil, err
	}
	return depts, nil
}

// GET ALL ROLES
func (r *UserRepository) FindAllRole() ([]models.Role, error) {
	var role []models.Role
	if err := database.DB.
		Find(&role).Error; err != nil {
		return nil, err
	}
	return role, nil
}

// READ ALL
func (r *UserRepository) FindAll() ([]models.User, error) {
	var users []models.User
	if err := database.DB.
		Preload("Department").
		Preload("Role").
		Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// READ ONE
func (r *UserRepository) FindByID(id uint64) (*models.User, error) {
	var user models.User
	err := r.DB.First(&user, id).Error
	return &user, err
}

// UPDATE
func (r *UserRepository) Update(user *models.User) error {
	return r.DB.Save(user).Error
}

// DELETE
func (r *UserRepository) Delete(id uint64) error {
	return r.DB.Delete(&models.User{}, id).Error
}

// PAGINATION FUNCTION
func (r *UserRepository) FindPaginated(page, pageSize int) (users []models.User, total int64, err error) {
	offset := (page - 1) * pageSize

	// total count
	if err = r.DB.Model(&models.User{}).Count(&total).Error; err != nil {
		return
	}

	err = r.DB.
		Preload("Department").
		Preload("Role").
		Limit(pageSize).
		Offset(offset).
		Order("id asc").
		Find(&users).Error

	return
}

// GetAssignableUsers returns users that can be assigned to tasks
func (ur *UserRepository) GetAssignableUsers() ([]models.User, error) {
	var users []models.User

	err := ur.DB.
		Where("is_active = ?", true).
		Order("first_name ASC").
		Find(&users).Error

	return users, err
}

func (r *UserRepository) GetUsersByProject(projectID uint64) ([]models.User, error) {
	var users []models.User
	err := r.DB.
		Joins("JOIN project_users ON project_users.user_id = users.id").
		Where("project_users.project_id = ?", projectID).
		Find(&users).Error
	return users, err
}

func (ur *UserRepository) GetByID(id uint64) (*models.User, error) {
	var user models.User
	if err := ur.DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
