package usecases

import (
	"time"

	"github.com/dmitrymomot/go-test-task-management-tool/app/domain"
)

// NewTaskInteractor is a factory function,
// returns an instance of TaskInteractor structure
func NewTaskInteractor(r domain.TaskRepository) *TaskInteractor {
	return &TaskInteractor{r}
}

// TaskInteractor structure
type TaskInteractor struct {
	r domain.TaskRepository
}

// GetAll returns the all tasks list
func (i *TaskInteractor) GetAll() ([]domain.Task, error) {
	tasks, err := i.r.GetAll()
	if err != nil {
		return []domain.Task{}, err
	}
	return tasks, nil
}

// GetNew returns the all new tasks list
func (i *TaskInteractor) GetNew() ([]domain.Task, error) {
	tasks, err := i.r.GetNew()
	if err != nil {
		return []domain.Task{}, err
	}
	return tasks, nil
}

// GetCompleted returns the all completed tasks list
func (i *TaskInteractor) GetCompleted() ([]domain.Task, error) {
	tasks, err := i.r.GetCompleted()
	if err != nil {
		return []domain.Task{}, err
	}
	return tasks, nil
}

// Store a new task
func (i *TaskInteractor) Store(title, description string) (domain.Task, error) {
	task := domain.NewTask(title, description)
	id, err := i.r.Store(task)
	if err != nil {
		return domain.Task{}, err
	}
	task, err = i.r.GetByID(id)
	if err != nil {
		return domain.Task{}, err
	}
	return task, nil
}

// Update a task
func (i *TaskInteractor) Update(id int64, title, description, status string) (domain.Task, error) {
	task, err := i.r.GetByID(id)
	if err != nil {
		return domain.Task{}, err
	}

	var updated = false

	if title != "" && task.Title != title {
		task.Title = title
		updated = true
	}
	if description != "" && task.Description != description {
		task.Description = description
		updated = true
	}
	if status != "" && (status == domain.CompletedTaskStatus || status == domain.NewTaskStatus) && status != task.Status {
		task.Status = status
		updated = true
		if status == domain.CompletedTaskStatus {
			now := time.Now()
			task.CompletedAt = &now
		} else if status == domain.NewTaskStatus {
			task.CompletedAt = nil
		}
	}

	if updated {
		err = i.r.Update(task)
		if err != nil {
			return task, err
		}
	}

	return task, nil
}

// Complete a task
func (i *TaskInteractor) Complete(id int64) (domain.Task, error) {
	task, err := i.r.GetByID(id)
	if err != nil {
		return task, err
	}
	task.Status = domain.CompletedTaskStatus
	now := time.Now()
	task.CompletedAt = &now

	err = i.r.Update(task)
	if err != nil {
		return task, err
	}

	return task, nil
}

// Delete a task
func (i *TaskInteractor) Delete(id int64) error {
	return i.r.Delete(id)
}
