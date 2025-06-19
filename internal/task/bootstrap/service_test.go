package bootstrap

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/yebrai/go-tasks-microservice/pkg/runner"
)

// MockConfig implementa runner.Config para testing
type MockConfig struct {
	serverAddress string
	mongoURI      string
	mongoDatabase string
	shouldFail    bool
}

func NewMockConfig() *MockConfig {
	return &MockConfig{
		serverAddress: ":8080",
		mongoURI:      "mongodb://localhost:27017/testdb",
		mongoDatabase: "testdb",
		shouldFail:    false,
	}
}

func (m *MockConfig) WithEmptyAddress() *MockConfig {
	m.serverAddress = ""
	return m
}

func (m *MockConfig) WithEmptyMongoURI() *MockConfig {
	m.mongoURI = ""
	return m
}

func (m *MockConfig) WithEmptyDatabase() *MockConfig {
	m.mongoDatabase = ""
	return m
}

func (m *MockConfig) Unmarshal(v interface{}) error {
	config := v.(*Config)
	config.Server.Address = m.serverAddress
	config.Mongo.URI = m.mongoURI
	config.Mongo.Database = m.mongoDatabase
	return nil
}

func TestService_Name(t *testing.T) {
	service := NewService()
	expected := "go-tasks-microservice"

	if service.Name() != expected {
		t.Errorf("Expected service name %s, got %s", expected, service.Name())
	}
}

func TestService_Configure_Success(t *testing.T) {
	service := NewService()
	mockConfig := NewMockConfig()

	err := service.configure(mockConfig)
	if err != nil {
		t.Errorf("Configuration should not fail: %v", err)
	}

	if service.config.Server.Address != ":8080" {
		t.Errorf("Expected server address :8080, got %s", service.config.Server.Address)
	}

	if service.config.Mongo.Database != "testdb" {
		t.Errorf("Expected database testdb, got %s", service.config.Mongo.Database)
	}

	if service.config.Mongo.URI != "mongodb://localhost:27017/testdb" {
		t.Errorf("Expected mongo URI mongodb://localhost:27017/testdb, got %s", service.config.Mongo.URI)
	}
}

func TestService_Configure_ValidationErrors(t *testing.T) {
	tests := []struct {
		name        string
		mockConfig  *MockConfig
		expectError bool
		errorMsg    string
	}{
		{
			name:        "empty server address",
			mockConfig:  NewMockConfig().WithEmptyAddress(),
			expectError: true,
			errorMsg:    "server address is required",
		},
		{
			name:        "empty mongo URI",
			mockConfig:  NewMockConfig().WithEmptyMongoURI(),
			expectError: true,
			errorMsg:    "mongo URI is required",
		},
		{
			name:        "empty mongo database",
			mockConfig:  NewMockConfig().WithEmptyDatabase(),
			expectError: true,
			errorMsg:    "mongo database name is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService()
			err := service.configure(tt.mockConfig)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestService_GracefulShutdown(t *testing.T) {
	service := NewService()
	service.config = &Config{
		Server: ServerConfig{Address: ":8081"},
	}

	// Simular un servidor HTTP simple para testing
	service.server = &http.Server{
		Addr: ":8081",
	}

	// Test graceful shutdown (sin iniciar el servidor)
	err := service.gracefulShutdown()
	if err != nil {
		t.Errorf("Graceful shutdown should not fail: %v", err)
	}
}

func TestService_Cleanup_WithNilProviders(t *testing.T) {
	service := NewService()
	service.providers = nil

	err := service.Cleanup()
	if err != nil {
		t.Errorf("Cleanup with nil providers should not fail: %v", err)
	}
}

// Benchmark para medir el tiempo de inicialización
func BenchmarkService_Configure(b *testing.B) {
	service := NewService()
	mockConfig := NewMockConfig()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := service.configure(mockConfig)
		if err != nil {
			b.Errorf("Configuration failed: %v", err)
		}
	}
}

// Test de integración básico (requiere MongoDB corriendo)
func TestService_Integration_Basic(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	service := NewService()
	mockConfig := NewMockConfig()

	// Solo testear configuración ya que para providers necesitaríamos MongoDB
	err := service.configure(mockConfig)
	if err != nil {
		t.Errorf("Integration test configuration failed: %v", err)
	}

	// Verificar que la configuración es válida
	if service.config == nil {
		t.Error("Config should not be nil after configuration")
	}
}
