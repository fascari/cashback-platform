package chain

import "context"

type (
	ID string

	MintTokenRequest struct {
		IdempotencyKey string
		WalletAddress  string
		TokenAmount    string
		ChainID        ID
	}

	MintTokenResult struct {
		Success         bool
		TransactionHash string
		BlockNumber     int64
		Status          string
		ErrorCode       string
		ErrorMessage    string
		Retryable       bool
	}

	BalanceResult struct {
		WalletAddress string
		Balance       string
		BlockNumber   int64
	}

	TransactionResult struct {
		TransactionHash string
		Status          string
		BlockNumber     int64
		Confirmations   int64
		GasUsed         int64
		Success         bool
	}

	Client interface {
		ChainID() ID
		MintToken(ctx context.Context, req MintTokenRequest) (*MintTokenResult, error)
		FetchBalance(ctx context.Context, walletAddress string) (*BalanceResult, error)
		FetchTransaction(ctx context.Context, txHash string) (*TransactionResult, error)
	}
)
