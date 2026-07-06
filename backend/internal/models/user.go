package models

import "time"

// User represents a user in the system.
// Password is stored as bcrypt hash — never in plaintext.
type User struct {
	ID           int64     `gorm:"primaryKey;column:id" json:"id"`
	Name         string    `gorm:"size:100;not null;column:name" json:"name"`
	Email        string    `gorm:"size:150;uniqueIndex;not null;column:email" json:"email"`
	PasswordHash string    `gorm:"size:255;not null;column:password_hash" json:"-"` // never exposed in API
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
}

func (User) TableName() string {
	return "users"
}
