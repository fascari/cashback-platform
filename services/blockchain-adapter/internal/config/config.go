package config

import "github.com/spf13/viper"

type (
	Config struct {
		App      AppConfig
		GRPC     GRPCConfig
		Database DatabaseConfig
		Ethereum EthereumConfig
		Wallet   WalletConfig
		Redis    RedisConfig
	}

	AppConfig struct {
		Name string
		Env  string
	}

	GRPCConfig struct {
		Port string
	}

	DatabaseConfig struct {
		DSN string
	}

	EthereumConfig struct {
		RPCURL          string
		ChainID         int64
		ContractAddress string
	}

	WalletConfig struct {
		Mnemonic       string
		DerivationPath string
	}

	RedisConfig struct {
		URL string
	}
)

func NewConfig() (*Config, error) {
	viper.AutomaticEnv()

	viper.SetDefault("APP_NAME", "blockchain-adapter")
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("GRPC_PORT", "50051")
	viper.SetDefault("POSTGRES_DSN_BLOCKCHAIN", "postgres://cashback_app:cashback_app@localhost:15432/blockchain_adapter_db?sslmode=disable&search_path=blockchain")
	viper.SetDefault("ETHEREUM_RPC_URL", "https://sepolia.infura.io/v3/YOUR_PROJECT_ID")
	viper.SetDefault("ETHEREUM_CHAIN_ID", 11155111) // Sepolia
	viper.SetDefault("WALLET_DERIVATION_PATH", "m/44'/60'/0'/0/0")
	viper.SetDefault("REDIS_URL", "redis://localhost:6379")

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
			DSN: viper.GetString("POSTGRES_DSN_BLOCKCHAIN"),
		},
		Ethereum: EthereumConfig{
			RPCURL:          viper.GetString("ETHEREUM_RPC_URL"),
			ChainID:         viper.GetInt64("ETHEREUM_CHAIN_ID"),
			ContractAddress: viper.GetString("CONTRACT_ADDRESS"),
		},
		Wallet: WalletConfig{
			Mnemonic:       viper.GetString("WALLET_MNEMONIC"),
			DerivationPath: viper.GetString("WALLET_DERIVATION_PATH"),
		},
		Redis: RedisConfig{
			URL: viper.GetString("REDIS_URL"),
		},
	}, nil
}
