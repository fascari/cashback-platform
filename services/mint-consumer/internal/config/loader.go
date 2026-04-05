package config

import (
	"github.com/cashback-platform/kit/logger"
	"github.com/spf13/viper"
)

// LoadDatabase reads database config from environment and panics on error.
func LoadDatabase() Database {
	return loadWithPanic(loadDatabase, "failed to load database config")
}

// LoadNATS reads NATS config from environment and panics on error.
func LoadNATS() NATS {
	return loadWithPanic(loadNATS, "failed to load NATS config")
}

// LoadGRPC reads gRPC config from environment and panics on error.
func LoadGRPC() GRPC {
	return loadWithPanic(loadGRPC, "failed to load gRPC config")
}

func loadDatabase() (Database, error) {
	viper.SetDefault("DATABASE_HOST", "localhost")
	viper.SetDefault("DATABASE_PORT", "5432")
	viper.SetDefault("DATABASE_USER", "postgres")
	viper.SetDefault("DATABASE_PASSWORD", "postgres")
	viper.SetDefault("DATABASE_NAME", "mint_consumer_db")
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

func loadNATS() (NATS, error) {
	viper.SetDefault("NATS_URL", "nats://localhost:4222")
	viper.AutomaticEnv()
	return NATS{URL: viper.GetString("NATS_URL")}, nil
}

func loadGRPC() (GRPC, error) {
	viper.SetDefault("BLOCKCHAIN_ADAPTER_GRPC_ADDRESS", "localhost:50051")
	viper.AutomaticEnv()
	return GRPC{BlockchainAdapterAddress: viper.GetString("BLOCKCHAIN_ADAPTER_GRPC_ADDRESS")}, nil
}

func loadWithPanic[T any](loader func() (T, error), msg string) T {
	v, err := loader()
	if err != nil {
		logger.Error(msg, "error", err)
		panic(err)
	}
	return v
}
