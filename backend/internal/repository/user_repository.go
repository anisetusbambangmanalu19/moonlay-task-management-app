package repository

import (
	"github.com/anisetusbambangmanalu19/moonlay-task-management-app/backend/internal/models"
	"gorm.io/gorm"
)

// UserRepository handles all database operations for the User model
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository instance
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// FindAll retrieves all users (excluding password hash)
func (r *UserRepository) FindAll() ([]models.User, error) {
	var users []models.User
	result := r.db.Select("id, name, email, created_at").Find(&users)
	return users, result.Error
}

// FindByEmail finds a user by email address (includes password hash for auth)
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// FindByID finds a user by their primary key ID
func (r *UserRepository) FindByID(id int64) (*models.User, error) {
	var user models.User
	result := r.db.Select("id, name, email, created_at").First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
