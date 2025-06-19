package task

import (
	"errors"
	"time"
)

// Status representa los posibles estados de una tarea
type Status string

const (
	StatusPending   Status = "pending"
	StatusCompleted Status = "completed"
	StatusCancelled Status = "cancelled"
)

// Errores de dominio
var (
	ErrTaskNotFound         = errors.New("task not found")
	ErrInvalidTaskID        = errors.New("invalid task ID")
	ErrInvalidTaskData      = errors.New("invalid task data")
	ErrTaskAlreadyCompleted = errors.New("task already completed")
	ErrTaskAlreadyCancelled = errors.New("task already cancelled")
)

// Task es la entidad principal del dominio
type Task struct {
	ID          string
	Title       string
	Description string
	Status      Status
	CreatedAt   time.Time
	DueDate     *time.Time
}

// NewTask Constructor para nuevas tareas
func NewTask(id, title, description string, dueDate *time.Time) (*Task, error) {
	if id == "" {
		return nil, ErrInvalidTaskID
	}

	if title == "" {
		return nil, ErrInvalidTaskData
	}

	return &Task{
		ID:          id,
		Title:       title,
		Description: description,
		Status:      StatusPending,
		CreatedAt:   time.Now(),
		DueDate:     dueDate,
	}, nil
}

// Complete marca una tarea como completada
func (t *Task) Complete() error {
	if t.Status == StatusCompleted {
		return ErrTaskAlreadyCompleted
	}

	if t.Status == StatusCancelled {
		return ErrTaskAlreadyCancelled
	}

	t.Status = StatusCompleted
	return nil
}

// Cancel marca una tarea como cancelada
func (t *Task) Cancel() error {
	if t.Status == StatusCompleted {
		return ErrTaskAlreadyCompleted
	}

	if t.Status == StatusCancelled {
		return ErrTaskAlreadyCancelled
	}

	t.Status = StatusCancelled
	return nil
}
