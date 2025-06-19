package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/yebrai/go-tasks-microservice/internal/task"
)

// TaskRepository implementación MongoDB del repositorio de tareas
type TaskRepository struct {
	collection *mongo.Collection
}

// NewTaskRepository crea una nueva instancia del repositorio
func NewTaskRepository(db *mongo.Database) *TaskRepository {
	return &TaskRepository{
		collection: db.Collection("tasks"),
	}
}

// TaskDocument representa la estructura de documento en MongoDB
type TaskDocument struct {
	ID          string     `bson:"_id"`
	Title       string     `bson:"title"`
	Description string     `bson:"description"`
	Status      string     `bson:"status"`
	CreatedAt   time.Time  `bson:"created_at"`
	DueDate     *time.Time `bson:"due_date,omitempty"`
}

// Save guarda una tarea en la base de datos
func (r *TaskRepository) Save(ctx context.Context, t *task.Task) error {
	doc := r.toDocument(t)

	_, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return fmt.Errorf("failed to save task: %w", err)
	}

	return nil
}

// FindByID busca una tarea por ID
func (r *TaskRepository) FindByID(ctx context.Context, id string) (*task.Task, error) {
	var doc TaskDocument

	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, task.ErrTaskNotFound
		}
		return nil, fmt.Errorf("failed to find task: %w", err)
	}

	return r.fromDocument(&doc), nil
}

// FindAll obtiene todas las tareas
func (r *TaskRepository) FindAll(ctx context.Context) ([]*task.Task, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to find tasks: %w", err)
	}
	defer cursor.Close(ctx)

	var tasks []*task.Task
	for cursor.Next(ctx) {
		var doc TaskDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, fmt.Errorf("failed to decode task: %w", err)
		}
		tasks = append(tasks, r.fromDocument(&doc))
	}

	return tasks, nil
}

// Update actualiza una tarea existente
func (r *TaskRepository) Update(ctx context.Context, t *task.Task) error {
	doc := r.toDocument(t)

	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": t.ID}, doc)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	return nil
}

// Delete elimina una tarea por ID
func (r *TaskRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}

// toDocument convierte una tarea del dominio a documento MongoDB
func (r *TaskRepository) toDocument(t *task.Task) *TaskDocument {
	return &TaskDocument{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Status:      string(t.Status),
		CreatedAt:   t.CreatedAt,
		DueDate:     t.DueDate,
	}
}

// fromDocument convierte un documento MongoDB a tarea del dominio
func (r *TaskRepository) fromDocument(doc *TaskDocument) *task.Task {
	return &task.Task{
		ID:          doc.ID,
		Title:       doc.Title,
		Description: doc.Description,
		Status:      task.Status(doc.Status),
		CreatedAt:   doc.CreatedAt,
		DueDate:     doc.DueDate,
	}
}
