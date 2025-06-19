package creator

import (
	"context"
	"fmt"

	"github.com/yebrai/go-tasks-microservice/internal/task"
	"github.com/yebrai/go-tasks-microservice/pkg/cqrs"
	"github.com/yebrai/go-tasks-microservice/pkg/id"
)

// CreateTaskCommandHandler maneja el comando de creación de tareas
type CreateTaskCommandHandler struct {
	repository  task.Repository
	idGenerator id.Generator
}

// NewCreateTaskCommandHandler crea una nueva instancia del handler
func NewCreateTaskCommandHandler(
	repository task.Repository,
	idGenerator id.Generator,
) *CreateTaskCommandHandler {
	return &CreateTaskCommandHandler{
		repository:  repository,
		idGenerator: idGenerator,
	}
}

// Handle procesa el comando de creación de tarea
func (h *CreateTaskCommandHandler) Handle(ctx context.Context, cmd cqrs.Command) error {
	createCmd, ok := cmd.(CreateTaskCommand)
	if !ok {
		return fmt.Errorf("invalid command type: expected CreateTaskCommand")
	}

	// Generar ID único para la tarea
	taskID := h.idGenerator.Generate()

	// Crear nueva tarea
	newTask, err := task.NewTask(taskID, createCmd.Title, createCmd.Description, createCmd.DueDate)
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	// Persistir la tarea
	if err := h.repository.Save(ctx, newTask); err != nil {
		return fmt.Errorf("failed to save task: %w", err)
	}

	// TODO: Aquí se podrían publicar eventos de dominio
	// eventBus.Publish(task.NewTaskCreatedEvent(newTask))

	return nil
}
