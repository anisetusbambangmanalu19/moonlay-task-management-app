package models

import "time"

// TaskStatus adalah enum untuk status task
type TaskStatus string

const (
	TaskStatusTodo       TaskStatus = "todo"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusDone       TaskStatus = "done"
)

// IsValid mengecek apakah nilai TaskStatus termasuk nilai yang diizinkan
func (s TaskStatus) IsValid() bool {
	return s == TaskStatusTodo || s == TaskStatusInProgress || s == TaskStatusDone
}

// Task merepresentasikan sebuah task di sistem.
// Menggunakan tipe ENUM PostgreSQL untuk field status.
type Task struct {
	ID          int64      `gorm:"primaryKey;column:id" json:"id"`
	Title       string     `gorm:"size:200;not null;column:title" json:"title"`
	Description string     `gorm:"type:text;column:description" json:"description"`
	Status      TaskStatus `gorm:"type:task_status;not null;default:'todo';column:status" json:"status"`
	Deadline    time.Time  `gorm:"not null;column:deadline" json:"deadline"`
	AssigneeID  int64      `gorm:"not null;column:assignee_id" json:"assignee_id"`
	CreatedBy   *int64     `gorm:"column:created_by" json:"created_by"`
	CreatedAt   time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at" json:"updated_at"`

	// Relasi (diisi lewat Preload)
	Assignee *User `gorm:"foreignKey:AssigneeID" json:"assignee,omitempty"`
	Creator  *User `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

func (Task) TableName() string {
	return "tasks"
}
