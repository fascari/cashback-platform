package config

import (
	"github.com/cashback-platform/services/cashback-service-api/pkg/logger"
	"github.com/spf13/viper"
)

type (
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
		SSLMode  string
	}

	NATS struct {
		URL string
	}

	GRPC struct {
		BlockchainAdapterAddress string
	}

	Server struct {
		Port string
	}
)

func LoadDatabase() Database {
	return loadConfigWithPanic(loadDatabaseConfig, "failed to load database config")
}

func LoadNATS() NATS {
	return loadConfigWithPanic(loadNATSConfig, "failed to load NATS config")
}

func LoadGRPC() GRPC {
	return loadConfigWithPanic(loadGRPCConfig, "failed to load GRPC config")
}

func LoadServer() Server {
	return loadConfigWithPanic(loadServerConfig, "failed to load server config")
}

func loadDatabaseConfig() (Database, error) {
	viper.SetDefault("DATABASE_HOST", "localhost")
	viper.SetDefault("DATABASE_PORT", "5432")
	viper.SetDefault("DATABASE_USER", "postgres")
	viper.SetDefault("DATABASE_PASSWORD", "postgres")
	viper.SetDefault("DATABASE_NAME", "cashback_service_db")
	viper.SetDefault("DATABASE_SSLMODE", "disable")
	viper.AutomaticEnv()

	return Database{
		Host:     viper.GetString("DATABASE_HOST"),
		Port:     viper.GetString("DATABASE_PORT"),
		User:     viper.GetString("DATABASE_USER"),
		Password: viper.GetString("DATABASE_PASSWORD"),
		Name:     viper.GetString("DATABASE_NAME"),
		SSLMode:  viper.GetString("DATABASE_SSLMODE"),
	}, nil
}

func loadNATSConfig() (NATS, error) {
	viper.SetDefault("NATS_URL", "nats://localhost:4222")
	viper.AutomaticEnv()
	return NATS{URL: viper.GetString("NATS_URL")}, nil
}

func loadGRPCConfig() (GRPC, error) {
	viper.SetDefault("BLOCKCHAIN_ADAPTER_GRPC_ADDRESS", "localhost:50051")
	viper.AutomaticEnv()
	return GRPC{BlockchainAdapterAddress: viper.GetString("BLOCKCHAIN_ADAPTER_GRPC_ADDRESS")}, nil
}

func loadServerConfig() (Server, error) {
	viper.SetDefault("SERVER_PORT", "8080")
	viper.AutomaticEnv()
	return Server{Port: viper.GetString("SERVER_PORT")}, nil
}

func loadConfigWithPanic[T any](loader func() (T, error), errorMsg string) T {
	config, err := loader()
	if err != nil {
		logger.Error(errorMsg, "error", err)
		panic(err)
	}
	return config
}
