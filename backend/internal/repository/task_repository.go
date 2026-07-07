package repository

import (
	"time"

	"github.com/anisetusbambangmanalu19/moonlay-task-management-app/backend/internal/models"
	"gorm.io/gorm"
)

// TaskContext adalah struct datar yang dipakai khusus untuk query konteks RAG chatbot.
// Menggunakan raw SQL — bukan model GORM — agar join dan filtering lebih fleksibel.
type TaskContext struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Status       string    `json:"status"`
	Deadline     time.Time `json:"deadline"`
	AssigneeName string    `json:"assignee_name"`
	CreatedAt    time.Time `json:"created_at"`
}

// TaskRepository menangani semua operasi database untuk model Task
type TaskRepository struct {
	db *gorm.DB
}

// NewTaskRepository membuat instance TaskRepository baru
func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

// FindAll mengambil semua task beserta informasi assignee-nya (via Preload)
func (r *TaskRepository) FindAll() ([]models.Task, error) {
	var tasks []models.Task
	result := r.db.Preload("Assignee").Order("created_at DESC").Find(&tasks)
	return tasks, result.Error
}

// FindByID mengambil satu task berdasarkan ID beserta semua relasinya
func (r *TaskRepository) FindByID(id int64) (*models.Task, error) {
	var task models.Task
	result := r.db.Preload("Assignee").Preload("Creator").First(&task, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &task, nil
}

// Create menyimpan task baru ke database
func (r *TaskRepository) Create(task *models.Task) error {
	return r.db.Create(task).Error
}

// Update menyimpan semua field dari task yang sudah ada
func (r *TaskRepository) Update(task *models.Task) error {
	// Gunakan Updates agar field bernilai nol tidak tertimpa secara tidak sengaja
	return r.db.Model(task).Updates(map[string]interface{}{
		"title":       task.Title,
		"description": task.Description,
		"status":      task.Status,
		"deadline":    task.Deadline,
		"assignee_id": task.AssigneeID,
	}).Error
}

// UpdateStatus hanya memperbarui field status dari task (untuk endpoint PATCH)
func (r *TaskRepository) UpdateStatus(id int64, status models.TaskStatus) error {
	result := r.db.Model(&models.Task{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Delete menghapus task berdasarkan ID
func (r *TaskRepository) Delete(id int64) error {
	result := r.db.Delete(&models.Task{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

	// GetTasksForChatbot memakai raw SQL untuk mengambil seluruh data task bagi konteks RAG chatbot.
	// Raw SQL sengaja dipakai di sini (sesuai spesifikasi), bukan GORM, agar lebih fleksibel.
func (r *TaskRepository) GetTasksForChatbot() ([]TaskContext, error) {
	var tasks []TaskContext

	// Join raw SQL — mengembalikan data konteks datar untuk prompt LLM
	result := r.db.Raw(`
		SELECT 
			t.id,
			t.title,
			COALESCE(t.description, '') AS description,
			t.status::text AS status,
			t.deadline,
			COALESCE(u.name, 'Tidak ada assignee') AS assignee_name,
			t.created_at
		FROM tasks t
		LEFT JOIN users u ON t.assignee_id = u.id
		ORDER BY t.deadline ASC
	`).Scan(&tasks)

	return tasks, result.Error
}
