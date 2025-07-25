package bootstrap

import (
	"context"
	"fmt"
	"net/http"
	"time"

	taskhttp "github.com/yebrai/go-tasks-microservice/internal/task/http"
	"github.com/yebrai/go-tasks-microservice/pkg/runner"
)

// Service implementa runner.Runner para el microservicio de tareas
type Service struct {
	config    *Config
	server    *http.Server
	providers *Providers
}

// NewService crea una nueva instancia del servicio
func NewService() *Service {
	return &Service{}
}

// Name retorna el nombre del microservicio
func (s *Service) Name() string {
	return "go-tasks-microservice"
}

// Run ejecuta el microservicio completo
func (s *Service) Run(ctx context.Context, config runner.Config) error {
	// 1. Cargar configuración
	if err := s.configure(config); err != nil {
		return fmt.Errorf("failed to configure service: %w", err)
	}

	// 2. Inicializar todas las dependencias
	providers, err := s.initProviders(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize providers: %w", err)
	}
	s.providers = providers

	// 3. Configurar servidor HTTP
	if err := s.setupHTTPServer(); err != nil {
		return fmt.Errorf("failed to setup HTTP server: %w", err)
	}

	// 4. Ejecutar servidor
	return s.serve(ctx)
}

// configure carga y valida la configuración del servicio
func (s *Service) configure(config runner.Config) error {
	s.config = &Config{}
	if err := config.Unmarshal(s.config); err != nil {
		return fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	// Validar configuración básica
	if s.config.Server.Address == "" {
		return fmt.Errorf("server address is required")
	}
	if s.config.Mongo.URI == "" {
		return fmt.Errorf("mongo URI is required")
	}
	if s.config.Mongo.Database == "" {
		return fmt.Errorf("mongo database name is required")
	}

	return nil
}

// initProviders inicializa todas las dependencias del microservicio
func (s *Service) initProviders(ctx context.Context) (*Providers, error) {
	providers, err := NewProviders(ctx, s.config)
	if err != nil {
		return nil, fmt.Errorf("failed to create providers: %w", err)
	}

	fmt.Printf("✅ Providers initialized successfully\n")
	fmt.Printf("   - MongoDB: %s\n", s.config.Mongo.Database)
	fmt.Printf("   - Command handlers: 1 registered\n")

	return providers, nil
}

// setupHTTPServer configura el servidor HTTP con todos los handlers
func (s *Service) setupHTTPServer() error {
	// Crear servidor HTTP con todas las dependencias inyectadas
	httpServer := taskhttp.NewServer(s.providers.CommandBus, s.providers.TaskRepository, s.providers.EventBus)

	// Configurar servidor HTTP con timeouts apropiados
	s.server = &http.Server{
		Addr:         s.config.Server.Address,
		Handler:      httpServer.Handler(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	fmt.Printf("✅ HTTP server configured on %s\n", s.config.Server.Address)
	return nil
}

// serve inicia el servidor HTTP y maneja el ciclo de vida
func (s *Service) serve(ctx context.Context) error {
	errChan := make(chan error, 1)

	// Iniciar servidor HTTP en goroutine
	go func() {
		fmt.Printf("🚀 HTTP server starting on %s\n", s.config.Server.Address)
		fmt.Printf("📋 Available endpoints:\n")
		fmt.Printf("   - GET  /health\n")
		fmt.Printf("   - GET  /ws/events (WebSocket)\n")
		fmt.Printf("   - GET  /api/v1/tasks\n")
		fmt.Printf("   - POST /api/v1/tasks\n")
		fmt.Printf("   - GET  /api/v1/tasks/:id\n")
		fmt.Printf("   - PUT  /api/v1/tasks/:id\n")

		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("HTTP server failed: %w", err)
		}
	}()

	// Esperar contexto cancelado o error del servidor
	select {
	case <-ctx.Done():
		fmt.Println("🛑 Received shutdown signal...")
		return s.gracefulShutdown()
	case err := <-errChan:
		return fmt.Errorf("server error: %w", err)
	}
}

// gracefulShutdown realiza un apagado controlado del servidor
func (s *Service) gracefulShutdown() error {
	fmt.Println("🔄 Starting graceful shutdown...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("failed to shutdown HTTP server: %w", err)
	}

	fmt.Println("✅ HTTP server stopped")
	return nil
}

// Cleanup limpia todos los recursos del servicio
func (s *Service) Cleanup() error {
	fmt.Println("🧹 Cleaning up resources...")

	if s.providers != nil {
		if err := s.providers.Cleanup(); err != nil {
			return fmt.Errorf("failed to cleanup providers: %w", err)
		}
	}

	fmt.Println("✅ Cleanup completed")
	return nil
}
