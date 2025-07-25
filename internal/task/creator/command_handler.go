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

// CompleteTaskCommandHandler maneja el comando para completar tareas
type CompleteTaskCommandHandler struct {
	repository task.Repository
	eventBus   events.EventBus
}

// NewCompleteTaskCommandHandler crea una nueva instancia del handler
func NewCompleteTaskCommandHandler(
	repository task.Repository,
	eventBus events.EventBus,
) *CompleteTaskCommandHandler {
	return &CompleteTaskCommandHandler{
		repository: repository,
		eventBus:   eventBus,
	}
}

// Handle maneja el comando CompleteTaskCommand
func (h *CompleteTaskCommandHandler) Handle(ctx context.Context, cmd cqrs.Command) error {
	completeCmd, ok := cmd.(CompleteTaskCommand)
	if !ok {
		return fmt.Errorf("invalid command type: expected CompleteTaskCommand")
	}

	// 1. Crear TaskID desde string
	taskID, err := task.NewID(completeCmd.ID)
	if err != nil {
		return fmt.Errorf("invalid task ID: %w", err)
	}

	// 2. Obtener la tarea
	existingTask, err := h.repository.FindByID(ctx, string(taskID))
	if err != nil {
		return fmt.Errorf("task not found: %w", err)
	}

	// 3. Completar la tarea (lógica de dominio)
	if err := existingTask.Complete(); err != nil {
		return fmt.Errorf("failed to complete task: %w", err)
	}

	// 4. Actualizar en el repositorio
	if err := h.repository.Update(ctx, existingTask); err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	// 5. Publicar evento
	event := task.NewTaskCompletedEvent(existingTask)
	if err := h.eventBus.Publish(ctx, event); err != nil {
		fmt.Printf("⚠️  Failed to publish task completed event: %v\n", err)
	}

	return nil
}
