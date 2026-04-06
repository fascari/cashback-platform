package domain

type (
	MintTokenRequest struct {
		IdempotencyKey string
		WalletAddress  string
		TokenAmount    string
		ChainID        string
	}

	MintResult struct {
		TransactionHash string
		BlockNumber     int64
		ErrorCode       string
		ErrorMessage    string
		Retryable       bool
	}
)
