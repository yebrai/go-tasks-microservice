package events

import (
	"context"
	"github.com/yebrai/go-tasks-microservice/internal/task"
)

type EventBus interface {
	Publish(ctx context.Context, event task.DomainEvent) error
	Close() error
}

type EventHandler interface {
	Handle(ctx context.Context, event task.DomainEvent) error
}
