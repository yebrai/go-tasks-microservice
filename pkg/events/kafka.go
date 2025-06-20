package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/yebrai/go-tasks-microservice/internal/task"
)

// KafkaEventBus implementación de EventBus usando Kafka
type KafkaEventBus struct {
	writer *kafka.Writer
}

// NewKafkaEventBus crea un EventBus que usa Kafka
func NewKafkaEventBus(writer *kafka.Writer) *KafkaEventBus {
	return &KafkaEventBus{
		writer: writer,
	}
}

// Publish publica un evento a Kafka
func (k *KafkaEventBus) Publish(ctx context.Context, event task.DomainEvent) error {
	// Serializar evento a JSON
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Crear mensaje Kafka
	message := kafka.Message{
		Key:   []byte(event.AggregateID()), // TaskID como partition key
		Value: payload,
		Time:  time.Now(),
		Headers: []kafka.Header{
			{Key: "event-type", Value: []byte(event.EventName())},
			{Key: "aggregate-id", Value: []byte(event.AggregateID())},
			{Key: "service", Value: []byte("go-tasks-microservice")},
			{Key: "version", Value: []byte("1.0")},
		},
	}

	// Publicar con timeout
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := k.writer.WriteMessages(ctx, message); err != nil {
		return fmt.Errorf("failed to publish event to kafka: %w", err)
	}

	fmt.Printf("📤 Event published to Kafka: %s - %s\n",
		event.EventName(),
		event.AggregateID())
	return nil
}

// Close no cierra el writer (se maneja en Providers.Cleanup)
func (k *KafkaEventBus) Close() error {
	// El KafkaWriter se cierra en Providers.Cleanup()
	fmt.Printf("✅ Kafka EventBus closed\n")
	return nil
}
