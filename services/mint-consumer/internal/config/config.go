package config

type (
	Database struct {
		DSN string
	}

	NATS struct {
		URL string
	}

	GRPC struct {
		BlockchainAdapterAddress string
	}
)
