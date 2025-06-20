package bootstrap

import (
	"context"
	"fmt"
	"github.com/yebrai/go-tasks-microservice/internal/task"
	"github.com/yebrai/go-tasks-microservice/internal/task/creator"
	taskmongo "github.com/yebrai/go-tasks-microservice/internal/task/mongo"
	"github.com/yebrai/go-tasks-microservice/pkg/cqrs"
	"github.com/yebrai/go-tasks-microservice/pkg/cqrs/inmem"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/yebrai/go-tasks-microservice/pkg/events"
	"github.com/yebrai/go-tasks-microservice/pkg/id"
	"github.com/yebrai/go-tasks-microservice/pkg/rabbitmq"
	"go.mongodb.org/mongo-driver/mongo"
)

type Providers struct {
	// INFRAESTRUCTURA BASE
	MongoClient *mongo.Client
	IDGenerator id.Generator

	// INFRAESTRUCTURA RabbitMQ
	RabbitMQClient *rabbitmq.Client

	// REPOSITORIOS (DOMAIN LAYER)
	TaskRepository task.Repository

	// CQRS (APPLICATION LAYER)
	CommandBus cqrs.CommandBus

	// EVENT SYSTEM (APPLICATION LAYER) - NUEVO
	EventBus events.EventBus

	// COMMAND HANDLERS (APPLICATION LAYER)
	CreateTaskHandler *creator.CreateTaskCommandHandler
}

// NewProviders crea e inicializa todas las dependencias del sistema
func NewProviders(ctx context.Context, config *Config) (*Providers, error) {
	providers := &Providers{}

	// INICIALIZACIÓN EN CAPAS (Bottom-Up)

	// 1. Infraestructura base (MongoDB, ID Generator)
	if err := providers.initInfrastructure(ctx, config); err != nil {
		return nil, fmt.Errorf("infrastructure initialization failed: %w", err)
	}

	// 2. Repositorios (Domain → Infrastructure adapters)
	if err := providers.initRepositories(config); err != nil {
		return nil, fmt.Errorf("repository initialization failed: %w", err)
	}

	// 3. CQRS buses
	if err := providers.initCQRS(); err != nil {
		return nil, fmt.Errorf("CQRS initialization failed: %w", err)
	}

	// 4. Event system
	if err := providers.initEventSystem(config); err != nil {
		return nil, fmt.Errorf("event system initialization failed: %w", err)
	}

	// 5. Command handlers
	if err := providers.initCommandHandlers(); err != nil {
		return nil, fmt.Errorf("command handlers initialization failed: %w", err)
	}

	return providers, nil
}

// FASE 1: INFRAESTRUCTURA BASE
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

// FASE 2: REPOSITORIOS
func (p *Providers) initRepositories(config *Config) error {
	database := p.MongoClient.Database(config.Mongo.Database)

	// Repository de tareas para operaciones CRUD
	p.TaskRepository = taskmongo.NewTaskRepository(database)

	fmt.Printf("✅ Repositories initialized\n")
	fmt.Printf("   - TaskRepository: MongoDB\n")

	return nil
}

// FASE 3: CQRS BUSES
func (p *Providers) initCQRS() error {
	p.CommandBus = inmem.NewCommandBus()

	fmt.Printf("✅ CQRS buses initialized\n")
	fmt.Printf("   - CommandBus: in-memory\n")

	return nil
}

// FASE 4: EVENT SYSTEM
func (p *Providers) initEventSystem(config *Config) error {
	if !config.RabbitMQ.Enabled {
		// RABBITMQ DESHABILITADO - usar NoOp EventBus
		p.EventBus = events.NewNoOpEventBus()
		fmt.Printf("⚠️  RabbitMQ disabled - using NoOp EventBus\n")
		return nil
	}

	// RABBITMQ HABILITADO - crear cliente real
	rabbitConfig := rabbitmq.ClientConfig{
		URL:      config.RabbitMQ.URL,
		Exchange: config.RabbitMQ.Exchange,
		Queue:    config.RabbitMQ.Queue,
	}

	// Crear y almacenar RabbitMQ Client
	client, err := rabbitmq.NewClient(rabbitConfig)
	if err != nil {
		return fmt.Errorf("failed to create RabbitMQ client: %w", err)
	}

	p.RabbitMQClient = client

	// Crear EventBus usando RabbitMQ
	p.EventBus = events.NewRabbitMQEventBus(client)

	fmt.Printf("✅ RabbitMQ EventBus initialized\n")
	fmt.Printf("   - URL: %s\n", config.RabbitMQ.URL)
	fmt.Printf("   - Exchange: %s\n", config.RabbitMQ.Exchange)

	return nil
}

// FASE 5: COMMAND HANDLERS
func (p *Providers) initCommandHandlers() error {
	// Handler para crear tareas CON EventBus inyectado
	p.CreateTaskHandler = creator.NewCreateTaskCommandHandler(
		p.TaskRepository,
		p.IDGenerator,
		p.EventBus,
	)

	// Registrar handlers en el command bus
	if err := p.CommandBus.Register(creator.CreateTaskCommandType, p.CreateTaskHandler); err != nil {
		return fmt.Errorf("failed to register CreateTaskCommandHandler: %w", err)
	}

	fmt.Printf("✅ Command handlers registered\n")
	fmt.Printf("   - CreateTaskCommand: ✓ (with EventBus)\n")

	return nil
}

// CLEANUP
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

	// Cleanup RabbitMQ Client
	if p.RabbitMQClient != nil {
		if err := p.RabbitMQClient.Close(); err != nil {
			errors = append(errors, fmt.Errorf("RabbitMQ client close error: %w", err))
		} else {
			fmt.Printf("✅ RabbitMQ client closed\n")
		}
	}

	// EventBus cleanup (principalmente para NoOp)
	if p.EventBus != nil {
		if err := p.EventBus.Close(); err != nil {
			errors = append(errors, fmt.Errorf("EventBus close error: %w", err))
		}
	}

	// Si hay errores, retornar el primero
	if len(errors) > 0 {
		return errors[0]
	}

	return nil
}
