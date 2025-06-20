package kafka

import (
	"crypto/tls"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

// ClientConfig configuración para cliente Kafka
type ClientConfig struct {
	Brokers     []string
	Topic       string
	Security    SecurityConfig
	Partitions  int
	Replication int
}

type SecurityConfig struct {
	Protocol string
	Username string
	Password string
}

// NewWriter crea un writer Kafka configurado
func NewWriter(config ClientConfig) *kafka.Writer {
	writerConfig := kafka.WriterConfig{
		Brokers:      config.Brokers,
		Topic:        config.Topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: 1,
		BatchSize:    100,
		Async:        false, // Síncrono para mejor error handling
	}

	// Configurar seguridad según protocolo
	switch config.Security.Protocol {
	case "PLAINTEXT":
		// Sin autenticación (desarrollo local)

	case "SASL_SSL":
		// Autenticación SSL (Upstash)
		mechanism := plain.Mechanism{
			Username: config.Security.Username,
			Password: config.Security.Password,
		}

		writerConfig.Dialer = &kafka.Dialer{
			SASLMechanism: mechanism,
			TLS:           &tls.Config{},
		}
	}

	return kafka.NewWriter(writerConfig)
}
