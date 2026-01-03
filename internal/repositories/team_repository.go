package repositories

import (
	"workms/internal/database"
	"workms/internal/models"

	"gorm.io/gorm"
)

type TeamRepository struct {
	DB *gorm.DB
}

func NewTeamRepository(db *gorm.DB) *TeamRepository {
	return &TeamRepository{DB: db}
}

// READ ALL
func (r *TeamRepository) FindAll() ([]models.Team, error) {
	var team []models.Team
	if err := database.DB.
		Find(&team).Error; err != nil {
		return nil, err
	}
	return team, nil
}

// CREATE
func (r *TeamRepository) Create(proj *models.Team) error {
	return r.DB.Create(proj).Error
}

// READ ONE
func (r *TeamRepository) FindByID(id uint64) (*models.Team, error) {
	var team models.Team
	err := r.DB.First(&team, id).Error
	return &team, err
}

// UPDATE
func (r *TeamRepository) Update(proj *models.Team) error {
	return r.DB.Save(proj).Error
}

// DELETE
func (r *TeamRepository) Delete(id uint64) error {
	return r.DB.Delete(&models.Team{}, id).Error
}

// PAGINATION FUNCTION
func (r *TeamRepository) FindPaginated(page, pageSize int) (team []models.Team, total int64, err error) {
	offset := (page - 1) * pageSize

	// total count
	if err = r.DB.Model(&models.Task{}).Count(&total).Error; err != nil {
		return
	}

	err = r.DB.
		Limit(pageSize).
		Offset(offset).
		Order("id asc").
		Find(&team).Error

	return
}
