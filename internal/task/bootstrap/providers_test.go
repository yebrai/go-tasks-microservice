package bootstrap

import (
	"testing"

	"github.com/yebrai/go-tasks-microservice/internal/task/creator"
	"github.com/yebrai/go-tasks-microservice/pkg/id"
)

func TestProviders_Structure(t *testing.T) {
	providers := &Providers{}

	// Verificar que la estructura esté bien definida
	if providers == nil {
		t.Error("Providers should be instantiable")
	}

	// Verificar campos principales
	if providers.MongoClient != nil {
		t.Error("MongoClient should be nil initially")
	}

	if providers.IDGenerator != nil {
		t.Error("IDGenerator should be nil initially")
	}

	if providers.TaskRepository != nil {
		t.Error("TaskRepository should be nil initially")
	}

	if providers.CommandBus != nil {
		t.Error("CommandBus should be nil initially")
	}

	if providers.CreateTaskHandler != nil {
		t.Error("CreateTaskHandler should be nil initially")
	}
}

func TestProviders_IDGenerator_Integration(t *testing.T) {
	// Test del generador de IDs de forma aislada
	generator := id.NewUniqueIDGenerator()

	if generator == nil {
		t.Error("ID generator should not be nil")
	}

	// Generar algunos IDs y verificar que son únicos
	ids := make(map[string]bool)
	for i := 0; i < 100; i++ {
		id := generator.Generate()

		if id == "" {
			t.Error("Generated ID should not be empty")
		}

		if ids[id] {
			t.Errorf("Duplicate ID generated: %s", id)
		}

		ids[id] = true
	}
}

func TestProviders_CommandHandlerType(t *testing.T) {
	// Verificar que el tipo de comando esté definido correctamente
	cmdType := creator.CreateTaskCommandType

	if cmdType == "" {
		t.Error("CreateTaskCommandType should not be empty")
	}

	expected := "task.command.create"
	if string(cmdType) != expected {
		t.Errorf("Expected command type %s, got %s", expected, string(cmdType))
	}
}

func TestProviders_Cleanup_Safety(t *testing.T) {
	providers := &Providers{}

	// Cleanup debería ser seguro incluso con providers vacíos
	err := providers.Cleanup()
	if err != nil {
		t.Errorf("Cleanup should be safe with empty providers: %v", err)
	}
}

// Test de configuración válida
func TestConfig_Valid(t *testing.T) {
	config := &Config{
		Server: ServerConfig{
			Address: ":8080",
		},
		Mongo: MongoConfig{
			URI:      "mongodb://localhost:27017/testdb",
			Database: "testdb",
		},
	}

	// Verificar configuración del servidor
	if config.Server.Address == "" {
		t.Error("Server address should not be empty")
	}

	if config.Server.Address != ":8080" {
		t.Errorf("Expected server address :8080, got %s", config.Server.Address)
	}

	// Verificar configuración de MongoDB
	if config.Mongo.URI == "" {
		t.Error("Mongo URI should not be empty")
	}

	if config.Mongo.Database == "" {
		t.Error("Mongo database should not be empty")
	}
}

// Benchmark de la estructura de providers
func BenchmarkProviders_Creation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		providers := &Providers{}
		_ = providers
	}
}
