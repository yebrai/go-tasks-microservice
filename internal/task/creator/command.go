package creator

import (
	"time"

	"github.com/yebrai/go-tasks-microservice/pkg/cqrs"
)

const CreateTaskCommandType cqrs.CommandType = "task.command.create"
const CompleteTaskCommandType cqrs.CommandType = "task.command.complete"

// CreateTaskCommand comando para crear una nueva tarea
type CreateTaskCommand struct {
	Title       string
	Description string
	DueDate     *time.Time
}

// Type implementa la interfaz Command
func (c CreateTaskCommand) Type() cqrs.CommandType {
	return CreateTaskCommandType
}

// CompleteTaskCommand comando para completar una tarea
type CompleteTaskCommand struct {
	ID string
}

// Type implementa la interfaz Command
func (c CompleteTaskCommand) Type() cqrs.CommandType {
	return CompleteTaskCommandType
}
