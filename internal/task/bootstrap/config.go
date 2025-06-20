package bootstrap

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Mongo    MongoConfig    `mapstructure:"mongo"`
	RabbitMQ RabbitMQConfig `mapstructure:"rabbitmq"`
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

type RabbitMQConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	URL      string `mapstructure:"url"`
	Exchange string `mapstructure:"exchange"`
	Queue    string `mapstructure:"queue"`
}
