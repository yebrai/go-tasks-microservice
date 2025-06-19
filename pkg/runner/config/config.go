package config

import (
	"fmt"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// Config implementa runner.Config usando koanf
type Config struct {
	k *koanf.Koanf
}

// Load carga la configuración desde archivo y variables de entorno
func Load() (*Config, error) {
	k := koanf.New(".")

	// Cargar desde archivo YAML si existe
	if err := k.Load(file.Provider("config.yaml"), yaml.Parser()); err != nil {
		// No es error crítico si no existe el archivo
		fmt.Printf("No config file found, using defaults and env vars\n")
	}

	// Cargar variables de entorno con prefijo
	if err := k.Load(env.Provider("TASK_", ".", func(s string) string {
		// Convertir TASK_SERVER_ADDRESS a server.address
		return strings.ToLower(strings.Replace(s[5:], "_", ".", -1))
	}), nil); err != nil {
		return nil, fmt.Errorf("failed to load environment variables: %w", err)
	}

	// Configurar valores por defecto
	setDefaults(k)

	return &Config{k: k}, nil
}

// Unmarshal deserializa la configuración en la estructura proporcionada
func (c *Config) Unmarshal(v interface{}) error {
	return c.k.Unmarshal("", v)
}

// setDefaults establece valores por defecto
func setDefaults(k *koanf.Koanf) {
	defaults := map[string]interface{}{
		"server.address": ":8080",
		"mongo.uri":      "mongodb://mongoroot:secret@localhost:27017/taskdb?authSource=admin",
		"mongo.database": "taskdb",
	}

	for key, value := range defaults {
		if !k.Exists(key) {
			k.Set(key, value)
		}
	}
}
