package cqrs

import "context"

// CommandType representa el tipo de comando
type CommandType string

// Command define el contrato para todos los comandos
type Command interface {
	Type() CommandType
}

// CommandHandler maneja la ejecución de comandos
type CommandHandler interface {
	Handle(ctx context.Context, cmd Command) error
}

// CommandBus maneja el dispatch de comandos
type CommandBus interface {
	Register(cmdType CommandType, handler CommandHandler) error
	Dispatch(ctx context.Context, cmd Command) error
}
