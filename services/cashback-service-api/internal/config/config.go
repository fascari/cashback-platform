package config

import (
	"github.com/spf13/viper"
)

type (
	Config struct {
		App      AppConfig
		Server   ServerConfig
		Database DatabaseConfig
		NATS     NATSConfig
		GRPC     GRPCConfig
	}

	AppConfig struct {
		Name string
		Env  string
	}

	ServerConfig struct {
		Port string
	}

	DatabaseConfig struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
		SSLMode  string
	}

	NATSConfig struct {
		URL string
	}

	GRPCConfig struct {
		BlockchainAdapterAddress string
	}
)

func NewConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// Environment variable bindings
	viper.SetEnvPrefix("")
	viper.AutomaticEnv()

	// Defaults
	viper.SetDefault("APP_NAME", "cashback-service-api")
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("DATABASE_HOST", "localhost")
	viper.SetDefault("DATABASE_PORT", "5432")
	viper.SetDefault("DATABASE_USER", "postgres")
	viper.SetDefault("DATABASE_PASSWORD", "postgres")
	viper.SetDefault("DATABASE_NAME", "cashback_service_db")
	viper.SetDefault("DATABASE_SSLMODE", "disable")
	viper.SetDefault("NATS_URL", "nats://localhost:4222")
	viper.SetDefault("BLOCKCHAIN_ADAPTER_GRPC_ADDRESS", "localhost:50051")

	// Try to read config file (optional)
	_ = viper.ReadInConfig()

	return &Config{
		App: AppConfig{
			Name: viper.GetString("APP_NAME"),
			Env:  viper.GetString("APP_ENV"),
		},
		Server: ServerConfig{
			Port: viper.GetString("SERVER_PORT"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("DATABASE_HOST"),
			Port:     viper.GetString("DATABASE_PORT"),
			User:     viper.GetString("DATABASE_USER"),
			Password: viper.GetString("DATABASE_PASSWORD"),
			Name:     viper.GetString("DATABASE_NAME"),
			SSLMode:  viper.GetString("DATABASE_SSLMODE"),
		},
		NATS: NATSConfig{
			URL: viper.GetString("NATS_URL"),
		},
		GRPC: GRPCConfig{
			BlockchainAdapterAddress: viper.GetString("BLOCKCHAIN_ADAPTER_GRPC_ADDRESS"),
		},
	}, nil
}
