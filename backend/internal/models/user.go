package models

import "time"

// User merepresentasikan user di sistem.
// Password disimpan sebagai hash bcrypt — tidak pernah plaintext.
type User struct {
	ID           int64     `gorm:"primaryKey;column:id" json:"id"`
	Name         string    `gorm:"size:100;not null;column:name" json:"name"`
	Email        string    `gorm:"size:150;uniqueIndex;not null;column:email" json:"email"`
	PasswordHash string    `gorm:"size:255;not null;column:password_hash" json:"-"` // tidak pernah diekspos di API
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
}

func (User) TableName() string {
	return "users"
}
