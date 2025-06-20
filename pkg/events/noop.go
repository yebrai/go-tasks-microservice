package events

import (
	"context"
	"fmt"
	"github.com/yebrai/go-tasks-microservice/internal/task"
)

// NoOpEventBus implementación que no hace nada (para desarrollo)
type NoOpEventBus struct{}

// NewNoOpEventBus crea un EventBus que solo hace logging
func NewNoOpEventBus() *NoOpEventBus {
	return &NoOpEventBus{}
}

// Publish solo hace logging, no envía eventos reales
func (n *NoOpEventBus) Publish(ctx context.Context, event task.DomainEvent) error {
	fmt.Printf("📝 [NoOp] Event: %s - AggregateID: %s - Time: %s\n",
		event.EventName(),
		event.AggregateID(),
		event.OccurredOn().Format("15:04:05"))
	return nil
}

// Close no hace nada
func (n *NoOpEventBus) Close() error {
	fmt.Printf("✅ NoOp EventBus closed\n")
	return nil
}
