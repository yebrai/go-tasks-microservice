package task

import (
	"time"
)

type DomainEvent interface {
	EventName() string
	AggregateID() string
	OccurredOn() time.Time
}

// BaseDomainEvent implementa la funcionalidad común de todos los eventos
type BaseDomainEvent struct {
	ID         string
	OccurredAt time.Time
}

type TaskCreatedEvent struct {
	BaseDomainEvent
	TaskID      string
	Title       string
	Description string
	DueDate     *time.Time
}

func NewTaskCreatedEvent(task *Task) *TaskCreatedEvent {
	return &TaskCreatedEvent{
		BaseDomainEvent: BaseDomainEvent{
			ID:         task.ID,
			OccurredAt: time.Now(),
		},
		TaskID:      task.ID,
		Title:       task.Title,
		Description: task.Description,
		DueDate:     task.DueDate,
	}
}

func (e TaskCreatedEvent) EventName() string {
	return "task.created"
}

func (e TaskCreatedEvent) AggregateID() string {
	return e.TaskID
}

func (e TaskCreatedEvent) OccurredOn() time.Time {
	return e.OccurredAt
}

type TaskCompletedEvent struct {
	BaseDomainEvent
	TaskID string
}

func NewTaskCompletedEvent(task *Task) *TaskCompletedEvent {
	return &TaskCompletedEvent{
		BaseDomainEvent: BaseDomainEvent{
			ID:         task.ID,
			OccurredAt: time.Now(),
		},
		TaskID: task.ID,
	}
}

func (e TaskCompletedEvent) EventName() string {
	return "task.completed"
}

func (e TaskCompletedEvent) AggregateID() string {
	return e.TaskID
}

func (e TaskCompletedEvent) OccurredOn() time.Time {
	return e.OccurredAt
}

type TaskCancelledEvent struct {
	BaseDomainEvent
	TaskID string
}

func NewTaskCancelledEvent(task *Task) *TaskCancelledEvent {
	return &TaskCancelledEvent{
		BaseDomainEvent: BaseDomainEvent{
			ID:         task.ID,
			OccurredAt: time.Now(),
		},
		TaskID: task.ID,
	}
}

func (e TaskCancelledEvent) EventName() string {
	return "task.cancelled"
}

func (e TaskCancelledEvent) AggregateID() string {
	return e.TaskID
}

func (e TaskCancelledEvent) OccurredOn() time.Time {
	return e.OccurredAt
}
