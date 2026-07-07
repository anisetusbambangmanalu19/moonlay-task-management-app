package repository

import (
	"github.com/anisetusbambangmanalu19/moonlay-task-management-app/backend/internal/models"
	"gorm.io/gorm"
)

// UserRepository menangani semua operasi database untuk model User
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository membuat instance UserRepository baru
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// FindAll mengambil semua user (tanpa password hash)
func (r *UserRepository) FindAll() ([]models.User, error) {
	var users []models.User
	result := r.db.Select("id, name, email, created_at").Find(&users)
	return users, result.Error
}

// FindByEmail mencari user berdasarkan alamat email (termasuk password hash untuk autentikasi)
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// FindByID mencari user berdasarkan primary key ID
func (r *UserRepository) FindByID(id int64) (*models.User, error) {
	var user models.User
	result := r.db.Select("id, name, email, created_at").First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
