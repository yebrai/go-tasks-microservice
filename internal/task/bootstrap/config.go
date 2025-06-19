package bootstrap

// Config contiene toda la configuración del microservicio
type Config struct {
	Server ServerConfig `mapstructure:"server"`
	Mongo  MongoConfig  `mapstructure:"mongo"`
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
