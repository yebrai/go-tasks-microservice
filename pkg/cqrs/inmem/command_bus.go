package inmem

import (
	"context"
	"fmt"
	"sync"

	"github.com/yebrai/go-tasks-microservice/pkg/cqrs"
)

// CommandBus implementación en memoria del bus de comandos
type CommandBus struct {
	handlers map[cqrs.CommandType]cqrs.CommandHandler
	mu       sync.RWMutex
}

// NewCommandBus crea una nueva instancia del bus de comandos
func NewCommandBus() *CommandBus {
	return &CommandBus{
		handlers: make(map[cqrs.CommandType]cqrs.CommandHandler),
	}
}

// Register registra un handler para un tipo de comando
func (b *CommandBus) Register(cmdType cqrs.CommandType, handler cqrs.CommandHandler) error {
	if cmdType == "" {
		return fmt.Errorf("command type cannot be empty")
	}
	if handler == nil {
		return fmt.Errorf("handler cannot be nil")
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	if _, exists := b.handlers[cmdType]; exists {
		return fmt.Errorf("handler already registered for command type: %s", cmdType)
	}

	b.handlers[cmdType] = handler
	return nil
}

// Dispatch ejecuta un comando usando el handler registrado
func (b *CommandBus) Dispatch(ctx context.Context, cmd cqrs.Command) error {
	if cmd == nil {
		return fmt.Errorf("command cannot be nil")
	}

	cmdType := cmd.Type()

	b.mu.RLock()
	handler, exists := b.handlers[cmdType]
	b.mu.RUnlock()

	if !exists {
		return fmt.Errorf("no handler registered for command type: %s", cmdType)
	}

	return handler.Handle(ctx, cmd)
}
