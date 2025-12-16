package config

import (
	"github.com/spf13/viper"
)

type (
	Config struct {
		App      AppConfig
		GRPC     GRPCConfig
		Database DatabaseConfig
	}

	AppConfig struct {
		Name string
		Env  string
	}

	GRPCConfig struct {
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
)

func NewConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	viper.SetEnvPrefix("")
	viper.AutomaticEnv()

	// Defaults
	viper.SetDefault("APP_NAME", "blockchain-adapter")
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("GRPC_PORT", "50051")
	viper.SetDefault("DATABASE_HOST", "localhost")
	viper.SetDefault("DATABASE_PORT", "5432")
	viper.SetDefault("DATABASE_USER", "postgres")
	viper.SetDefault("DATABASE_PASSWORD", "postgres")
	viper.SetDefault("DATABASE_NAME", "blockchain_adapter_db")
	viper.SetDefault("DATABASE_SSLMODE", "disable")

	_ = viper.ReadInConfig()

	return &Config{
		App: AppConfig{
			Name: viper.GetString("APP_NAME"),
			Env:  viper.GetString("APP_ENV"),
		},
		GRPC: GRPCConfig{
			Port: viper.GetString("GRPC_PORT"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("DATABASE_HOST"),
			Port:     viper.GetString("DATABASE_PORT"),
			User:     viper.GetString("DATABASE_USER"),
			Password: viper.GetString("DATABASE_PASSWORD"),
			Name:     viper.GetString("DATABASE_NAME"),
			SSLMode:  viper.GetString("DATABASE_SSLMODE"),
		},
	}, nil
}
