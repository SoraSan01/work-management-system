package repositories

import (
	"workms/internal/models"

	"gorm.io/gorm"
)

type RoleRepository struct {
	DB *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{DB: db}
}

// CREATE
func (r *RoleRepository) Create(role *models.Role) error {
	return r.DB.Create(role).Error
}

// READ ALL
func (r *RoleRepository) FindAll() ([]models.Role, error) {
	var role []models.Role
	err := r.DB.Find(&role).Error
	return role, err
}

// READ ONE
func (r *RoleRepository) FindByID(id uint64) (*models.Role, error) {
	var role models.Role
	err := r.DB.First(&role, id).Error
	return &role, err
}

// UPDATE
func (r *RoleRepository) Update(role *models.Role) error {
	return r.DB.Save(role).Error
}

// DELETE
func (r *RoleRepository) Delete(id uint64) error {
	return r.DB.Delete(&models.Role{}, id).Error
}
