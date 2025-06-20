package creator

import (
	"context"
	"fmt"

	"github.com/yebrai/go-tasks-microservice/internal/task"
	"github.com/yebrai/go-tasks-microservice/pkg/cqrs"
	"github.com/yebrai/go-tasks-microservice/pkg/events"
	"github.com/yebrai/go-tasks-microservice/pkg/id"
)

// ACTUALIZAR STRUCT para incluir EventBus
type CreateTaskCommandHandler struct {
	repository  task.Repository
	idGenerator id.Generator
	eventBus    events.EventBus
}

// ACTUALIZAR CONSTRUCTOR
func NewCreateTaskCommandHandler(
	repository task.Repository,
	idGenerator id.Generator,
	eventBus events.EventBus,
) *CreateTaskCommandHandler {
	return &CreateTaskCommandHandler{
		repository:  repository,
		idGenerator: idGenerator,
		eventBus:    eventBus, // ASIGNAR
	}
}

// ACTUALIZAR Handle para publicar eventos
func (h *CreateTaskCommandHandler) Handle(ctx context.Context, cmd cqrs.Command) error {
	createCmd, ok := cmd.(CreateTaskCommand)
	if !ok {
		return fmt.Errorf("invalid command type: expected CreateTaskCommand")
	}

	// 1. Generar ID único para la tarea
	taskID := h.idGenerator.Generate()

	// 2. Crear nueva tarea (lógica de dominio)
	newTask, err := task.NewTask(taskID, createCmd.Title, createCmd.Description, createCmd.DueDate)
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	// 3. Persistir la tarea (operación principal)
	if err := h.repository.Save(ctx, newTask); err != nil {
		return fmt.Errorf("failed to save task: %w", err)
	}

	// 4. ✅ PUBLICAR EVENTO (NUEVO)
	event := task.NewTaskCreatedEvent(newTask)
	if err := h.eventBus.Publish(ctx, event); err != nil {
		fmt.Printf("⚠️  Failed to publish event: %v\n", err)
		// - return fmt.Errorf("failed to publish event: %w", err)
	}

	return nil
}
