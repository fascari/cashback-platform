package config

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
)
