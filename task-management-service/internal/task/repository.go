package task

import (
	"context"
)

// Repository define las operaciones disponibles para persistir tareas
type Repository interface {
	Save(ctx context.Context, task *Task) error
	FindByID(ctx context.Context, id string) (*Task, error)
	FindAll(ctx context.Context) ([]*Task, error)
	Update(ctx context.Context, task *Task) error
	Delete(ctx context.Context, id string) error
}
