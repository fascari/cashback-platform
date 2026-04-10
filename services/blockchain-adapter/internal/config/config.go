package config

import "github.com/spf13/viper"

const (
	envAppName        = "APP_NAME"
	envAppEnv         = "APP_ENV"
	envGRPCPort       = "GRPC_PORT"
	envPostgresDSN    = "POSTGRES_DSN_BLOCKCHAIN"
	envEthereumRPCURL = "ETHEREUM_RPC_URL"
	envEthereumChain  = "ETHEREUM_CHAIN_ID"
	envContractAddr   = "CONTRACT_ADDRESS"
	envWalletMnemonic = "WALLET_MNEMONIC"
	envDerivationPath = "WALLET_DERIVATION_PATH"
	envRedisURL       = "REDIS_URL"

	defaultAppName        = "blockchain-adapter"
	defaultAppEnv         = "development"
	defaultGRPCPort       = "50051"
	defaultPostgresDSN    = "postgres://cashback_app:cashback_app@localhost:15432/blockchain_adapter_db?sslmode=disable&search_path=blockchain"
	defaultEthereumRPCURL = "https://sepolia.infura.io/v3/YOUR_PROJECT_ID"
	defaultSepoliaChainID = 11155111
	defaultDerivationPath = "m/44'/60'/0'/0/0"
	defaultRedisURL       = "redis://localhost:6379"
)

func NewConfig() (*Config, error) {
	viper.AutomaticEnv()

	viper.SetDefault(envAppName, defaultAppName)
	viper.SetDefault(envAppEnv, defaultAppEnv)
	viper.SetDefault(envGRPCPort, defaultGRPCPort)
	viper.SetDefault(envPostgresDSN, defaultPostgresDSN)
	viper.SetDefault(envEthereumRPCURL, defaultEthereumRPCURL)
	viper.SetDefault(envEthereumChain, defaultSepoliaChainID)
	viper.SetDefault(envDerivationPath, defaultDerivationPath)
	viper.SetDefault(envRedisURL, defaultRedisURL)

	_ = viper.ReadInConfig()

	return &Config{
		App: AppConfig{
			Name: viper.GetString(envAppName),
			Env:  viper.GetString(envAppEnv),
		},
		GRPC: GRPCConfig{
			Port: viper.GetString(envGRPCPort),
		},
		Database: DatabaseConfig{
			DSN: viper.GetString(envPostgresDSN),
		},
		Ethereum: EthereumConfig{
			RPCURL:          viper.GetString(envEthereumRPCURL),
			ChainID:         viper.GetInt64(envEthereumChain),
			ContractAddress: viper.GetString(envContractAddr),
		},
		Wallet: WalletConfig{
			Mnemonic:       viper.GetString(envWalletMnemonic),
			DerivationPath: viper.GetString(envDerivationPath),
		},
		Redis: RedisConfig{
			URL: viper.GetString(envRedisURL),
		},
	}, nil
}
