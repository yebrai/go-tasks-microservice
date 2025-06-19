package runner

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	runnerconfig "github.com/yebrai/go-tasks-microservice/pkg/runner/config"
)

// RunnerFactory función factory para crear runners
type RunnerFactory func(ctx context.Context) Runner

// Run ejecuta un runner con manejo de señales y configuración automática
func Run(factory RunnerFactory) {
	// Crear contexto con cancelación por señales
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Cargar configuración
	config, err := runnerconfig.Load()
	if err != nil {
		fmt.Printf("❌ Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Crear runner usando factory
	runner := factory(ctx)

	fmt.Printf("🚀 Starting %s...\n", runner.Name())

	// Ejecutar runner principal
	errChan := make(chan error, 1)
	go func() {
		if err := runner.Run(ctx, config); err != nil {
			errChan <- fmt.Errorf("runner execution failed: %w", err)
		}
	}()

	// Esperar terminación o error
	select {
	case <-ctx.Done():
		fmt.Printf("\n🛑 Received shutdown signal for %s\n", runner.Name())
	case err := <-errChan:
		fmt.Printf("❌ %s failed: %v\n", runner.Name(), err)
		cancel()
	}

	// Cleanup de recursos
	if err := runner.Cleanup(); err != nil {
		fmt.Printf("⚠️  Cleanup warning: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ %s stopped gracefully\n", runner.Name())
}
