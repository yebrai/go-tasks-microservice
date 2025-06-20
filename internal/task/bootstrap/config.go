package bootstrap

type Config struct {
	Server ServerConfig `mapstructure:"server"`
	Mongo  MongoConfig  `mapstructure:"mongo"`
	Kafka  KafkaConfig  `mapstructure:"kafka"`
}

// ServerConfig configuración del servidor HTTP
type ServerConfig struct {
	Address string `mapstructure:"address"`
}

// MongoConfig configuración de la base de datos MongoDB
type MongoConfig struct {
	URI      string `mapstructure:"uri"`
	Database string `mapstructure:"database"`
}

type KafkaConfig struct {
	Enabled     bool           `mapstructure:"enabled"`
	Brokers     []string       `mapstructure:"brokers"`
	Topic       string         `mapstructure:"topic"`
	Security    SecurityConfig `mapstructure:"security"`
	Partitions  int            `mapstructure:"partitions"`
	Replication int            `mapstructure:"replication"`
}

type SecurityConfig struct {
	Protocol string `mapstructure:"protocol"` // PLAINTEXT, SASL_SSL
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}
