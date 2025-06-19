package bootstrap

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/yebrai/go-tasks-microservice/internal/task"
	"github.com/yebrai/go-tasks-microservice/internal/task/creator"
	taskmongo "github.com/yebrai/go-tasks-microservice/internal/task/mongo"
	"github.com/yebrai/go-tasks-microservice/pkg/cqrs"
	"github.com/yebrai/go-tasks-microservice/pkg/cqrs/inmem"
	"github.com/yebrai/go-tasks-microservice/pkg/id"
)

// Providers contiene todas las dependencias inyectadas del microservicio
type Providers struct {
	// Infrastructure layer
	MongoClient *mongo.Client
	IDGenerator id.Generator

	// Domain layer - Repositories
	TaskRepository task.Repository

	// Application layer - CQRS
	CommandBus cqrs.CommandBus

	// Application layer - Command Handlers
	CreateTaskHandler *creator.CreateTaskCommandHandler
}

// NewProviders crea e inicializa todas las dependencias del sistema
func NewProviders(ctx context.Context, config *Config) (*Providers, error) {
	providers := &Providers{}

	// 1. Inicializar infraestructura base
	if err := providers.initInfrastructure(ctx, config); err != nil {
		return nil, fmt.Errorf("infrastructure initialization failed: %w", err)
	}

	// 2. Inicializar capa de persistencia
	if err := providers.initRepositories(config); err != nil {
		return nil, fmt.Errorf("repository initialization failed: %w", err)
	}

	// 3. Inicializar CQRS buses
	if err := providers.initCQRS(); err != nil {
		return nil, fmt.Errorf("CQRS initialization failed: %w", err)
	}

	// 4. Inicializar y registrar command handlers
	if err := providers.initCommandHandlers(); err != nil {
		return nil, fmt.Errorf("command handlers initialization failed: %w", err)
	}

	return providers, nil
}

// initInfrastructure inicializa la infraestructura base (MongoDB, generadores, etc.)
func (p *Providers) initInfrastructure(ctx context.Context, config *Config) error {
	// Conexión a MongoDB con configuración robusta
	clientOptions := options.Client().
		ApplyURI(config.Mongo.URI).
		SetMaxPoolSize(10).
		SetMinPoolSize(2)

	mongoClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Verificar conexión con ping
	if err := mongoClient.Ping(ctx, nil); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	p.MongoClient = mongoClient

	// Inicializar generador de IDs único
	p.IDGenerator = id.NewUniqueIDGenerator()

	fmt.Printf("✅ Infrastructure initialized\n")
	fmt.Printf("   - MongoDB connected: %s\n", config.Mongo.Database)
	fmt.Printf("   - ID generator: UUID\n")

	return nil
}

// initRepositories inicializa todos los repositorios de persistencia
func (p *Providers) initRepositories(config *Config) error {
	database := p.MongoClient.Database(config.Mongo.Database)

	// Repository de tareas para operaciones CRUD
	p.TaskRepository = taskmongo.NewTaskRepository(database)

	fmt.Printf("✅ Repositories initialized\n")
	fmt.Printf("   - TaskRepository: MongoDB\n")

	return nil
}

// initCQRS inicializa los buses de comandos y consultas
func (p *Providers) initCQRS() error {
	// Command Bus en memoria (escalable a Redis/RabbitMQ después)
	p.CommandBus = inmem.NewCommandBus()

	fmt.Printf("✅ CQRS buses initialized\n")
	fmt.Printf("   - CommandBus: in-memory\n")

	return nil
}

// initCommandHandlers inicializa y registra todos los command handlers
func (p *Providers) initCommandHandlers() error {
	// Handler para crear tareas
	p.CreateTaskHandler = creator.NewCreateTaskCommandHandler(
		p.TaskRepository,
		p.IDGenerator,
	)

	// Registrar handlers en el command bus
	if err := p.CommandBus.Register(creator.CreateTaskCommandType, p.CreateTaskHandler); err != nil {
		return fmt.Errorf("failed to register CreateTaskCommandHandler: %w", err)
	}

	fmt.Printf("✅ Command handlers registered\n")
	fmt.Printf("   - CreateTaskCommand: ✓\n")

	return nil
}

// Cleanup limpia todos los recursos y conexiones
func (p *Providers) Cleanup() error {
	var errors []error

	// Cerrar conexión MongoDB
	if p.MongoClient != nil {
		if err := p.MongoClient.Disconnect(context.Background()); err != nil {
			errors = append(errors, fmt.Errorf("MongoDB disconnect error: %w", err))
		} else {
			fmt.Printf("✅ MongoDB connection closed\n")
		}
	}

	// Si hay errores, retornar el primero
	if len(errors) > 0 {
		return errors[0]
	}

	return nil
}
