package domain

import "time"

// Predefined task statuses
const (
	NewTaskStatus       = "new"
	CompletedTaskStatus = "done"
)

// NewTask is factory function,
// returns a new instance od Task structure
func NewTask(title, description string) Task {
	return Task{
		Title:       title,
		Description: description,
		Status:      NewTaskStatus,
		CreatedAt:   time.Now(),
	}
}

// Task item structure
type Task struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// TaskRepository interface
type TaskRepository interface {
	GetByID(id int64) (Task, error)
	GetAll() ([]Task, error)
	GetNew() ([]Task, error)
	GetCompleted() ([]Task, error)
	Store(Task) (int64, error)
	Update(Task) error
	Delete(id int64) error
}
