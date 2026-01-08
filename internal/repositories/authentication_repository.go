package repositories

import (
	"errors"
	"workms/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthenticationRepository struct {
	DB *gorm.DB
}

func NewAuthenticationRepository(db *gorm.DB) *AuthenticationRepository {
	return &AuthenticationRepository{DB: db}
}

// Verify user email & password
func (r *AuthenticationRepository) Login(email, password string) (*models.User, error) {
	var user models.User
	if err := r.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	// Compare password hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	return &user, nil
}

func (ar *AuthenticationRepository) Authenticate(email, password string) (*models.User, error) {
	var user models.User

	if err := ar.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Compare hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid password")
	}

	return &user, nil
}
