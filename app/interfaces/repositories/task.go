package repositories

import (
	"database/sql"
	"errors"

	"github.com/dmitrymomot/go-test-task-management-tool/app/domain"
)

// Predefined errors
var (
	ErrNotFound = errors.New("task not found")
)

// NewTaskRepository is a factory function,
// returns an instance of TaskRepository structure
func NewTaskRepository(db DbHandler) *TaskRepository {
	return &TaskRepository{db}
}

// TaskRepository structure
type TaskRepository struct {
	db DbHandler
}

// GetByID retrieves task from storage by id
func (r *TaskRepository) GetByID(id int64) (domain.Task, error) {
	q := `SELECT id, title, description, status, created_at, completed_at FROM tasks WHERE id=? LIMIT 1`
	task := domain.Task{}
	err := r.db.QueryRow(q, id).Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.CompletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Task{}, ErrNotFound
		}
		return domain.Task{}, err
	}
	return task, nil
}

// GetAll returns all tasks from storage
func (r *TaskRepository) GetAll() ([]domain.Task, error) {
	q := `SELECT id, title, description, status, created_at, completed_at FROM tasks ORDER BY created_at DESC`
	rows, err := r.db.Query(q)
	if err != nil {
		return []domain.Task{}, err
	}

	tasks := []domain.Task{}
	for rows.Next() {
		task := domain.Task{}
		err = rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.CompletedAt)
		if err != nil {
			return []domain.Task{}, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// GetNew returns all new tasks from storage
func (r *TaskRepository) GetNew() ([]domain.Task, error) {
	return r.getByStatus(domain.NewTaskStatus)
}

// GetCompleted returns all completed tasks from storage
func (r *TaskRepository) GetCompleted() ([]domain.Task, error) {
	return r.getByStatus(domain.CompletedTaskStatus)
}

// getByStatus returns all tasks by status
func (r *TaskRepository) getByStatus(status string) ([]domain.Task, error) {
	q := `SELECT id, title, description, status, created_at, completed_at
		FROM tasks
		WHERE status=?
		ORDER BY created_at DESC`
	rows, err := r.db.Query(q, status)
	if err != nil {
		return []domain.Task{}, err
	}

	tasks := []domain.Task{}
	for rows.Next() {
		task := domain.Task{}
		err = rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.CompletedAt)
		if err != nil {
			return []domain.Task{}, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// Store new task into storage
func (r *TaskRepository) Store(task domain.Task) (id int64, err error) {
	q := `INSERT INTO tasks (title, description, status, created_at) VALUES (?, ?, ?, ?)`
	res, err := r.db.Execute(q, task.Title, task.Description, task.Status, task.CreatedAt)
	if err != nil {
		return id, err
	}
	id, err = res.LastInsertId()
	if err != nil {
		return id, err
	}
	return id, nil
}

// Update task
func (r *TaskRepository) Update(task domain.Task) error {
	q := `UPDATE tasks SET title=?, description=?, status=?, completed_at=? WHERE id=?`
	_, err := r.db.Execute(q, task.Title, task.Description, task.Status, task.CompletedAt, task.ID)
	if err != nil {
		return err
	}
	return nil
}

// Delete task from storage
func (r *TaskRepository) Delete(id int64) error {
	q := `DELETE FROM tasks WHERE id=?`
	_, err := r.db.Execute(q, id)
	if err != nil {
		return err
	}
	return nil
}
