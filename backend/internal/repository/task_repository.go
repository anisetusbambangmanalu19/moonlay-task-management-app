package repository

import (
	"time"

	"github.com/anisetusbambangmanalu19/moonlay-task-management-app/backend/internal/models"
	"gorm.io/gorm"
)

// TaskContext is a flat struct used specifically for the chatbot RAG context query.
// Uses raw SQL — not GORM model — to allow flexible joins and filtering.
type TaskContext struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Status       string    `json:"status"`
	Deadline     time.Time `json:"deadline"`
	AssigneeName string    `json:"assignee_name"`
	CreatedAt    time.Time `json:"created_at"`
}

// TaskRepository handles all database operations for the Task model
type TaskRepository struct {
	db *gorm.DB
}

// NewTaskRepository creates a new TaskRepository instance
func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

// FindAll retrieves all tasks with their assignee info (via Preload)
func (r *TaskRepository) FindAll() ([]models.Task, error) {
	var tasks []models.Task
	result := r.db.Preload("Assignee").Order("created_at DESC").Find(&tasks)
	return tasks, result.Error
}

// FindByID retrieves a single task by ID with all associations
func (r *TaskRepository) FindByID(id int64) (*models.Task, error) {
	var task models.Task
	result := r.db.Preload("Assignee").Preload("Creator").First(&task, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &task, nil
}

// Create inserts a new task into the database
func (r *TaskRepository) Create(task *models.Task) error {
	return r.db.Create(task).Error
}

// Update saves all fields of an existing task
func (r *TaskRepository) Update(task *models.Task) error {
	// Use Updates to avoid overwriting zero-value fields accidentally
	return r.db.Model(task).Updates(map[string]interface{}{
		"title":       task.Title,
		"description": task.Description,
		"status":      task.Status,
		"deadline":    task.Deadline,
		"assignee_id": task.AssigneeID,
	}).Error
}

// UpdateStatus updates only the status field of a task (for PATCH endpoint)
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

// Delete soft-deletes a task by ID
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

// GetTasksForChatbot uses raw SQL to fetch all task data for the chatbot RAG context.
// Raw SQL is intentionally used here (as per spec) instead of GORM for flexibility.
func (r *TaskRepository) GetTasksForChatbot() ([]TaskContext, error) {
	var tasks []TaskContext

	// Raw SQL join — returns flat context data for the LLM prompt
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
