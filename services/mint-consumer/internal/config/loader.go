package config

import (
	"github.com/cashback-platform/kit/logger"
	"github.com/spf13/viper"
)

func LoadDatabase() Database {
	return loadWithPanic(readDatabase, "failed to load database config")
}

func LoadNATS() NATS {
	return loadWithPanic(readNATS, "failed to load NATS config")
}

func LoadGRPC() GRPC {
	return loadWithPanic(readGRPC, "failed to load gRPC config")
}

func readDatabase() (Database, error) {
	viper.SetDefault("POSTGRES_DSN_MINT", "postgres://cashback_app:cashback_app@localhost:15432/mint_consumer_db?sslmode=disable&search_path=mint")
	viper.AutomaticEnv()
	return Database{
		DSN: viper.GetString("POSTGRES_DSN_MINT"),
	}, nil
}

func readNATS() (NATS, error) {
	viper.SetDefault("NATS_URL", "nats://localhost:4222")
	viper.AutomaticEnv()
	return NATS{URL: viper.GetString("NATS_URL")}, nil
}

func readGRPC() (GRPC, error) {
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
