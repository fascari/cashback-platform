package config

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
