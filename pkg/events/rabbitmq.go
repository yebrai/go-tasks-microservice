package events

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/yebrai/go-tasks-microservice/internal/task"
	"github.com/yebrai/go-tasks-microservice/pkg/rabbitmq"
)

type RabbitMQEventBus struct {
	client *rabbitmq.Client
}

func NewRabbitMQEventBus(client *rabbitmq.Client) *RabbitMQEventBus {
	return &RabbitMQEventBus{
		client: client,
	}
}

func (r *RabbitMQEventBus) Publish(_ context.Context, event task.DomainEvent) error {
	// Serializar evento a JSON
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Usar event name como routing key
	routingKey := event.EventName()

	// Publicar a RabbitMQ con context handling
	if err := r.client.Publish(routingKey, payload); err != nil {
		return fmt.Errorf("failed to publish event to rabbitmq: %w", err)
	}

	fmt.Printf("📤 Event published to RabbitMQ: %s - %s\n",
		event.EventName(),
		event.AggregateID())
	return nil
}

func (r *RabbitMQEventBus) Close() error {
	fmt.Printf("✅ RabbitMQ EventBus closed\n")
	return nil
}
