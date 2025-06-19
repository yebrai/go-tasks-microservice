package creator

import (
	"time"

	"github.com/yebrai/go-tasks-microservice/pkg/cqrs"
)

const CreateTaskCommandType cqrs.CommandType = "task.command.create"

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
