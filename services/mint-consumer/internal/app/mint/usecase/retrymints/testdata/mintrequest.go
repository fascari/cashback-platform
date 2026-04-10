package testdata

import (
	"github.com/cashback-platform/services/mint-consumer/internal/app/mint/domain"
	"github.com/google/uuid"
)

const (
	MintRequestID int64  = 42
	WalletAddress        = "0xABCDEF1234567890"
	TokenAmount          = "100"
	RetryCount           = 2
)

func FailedRetryableMintRequest() domain.MintRequest {
	return domain.MintRequest{
		ID:             MintRequestID,
		WalletAddress:  WalletAddress,
		TokenAmount:    TokenAmount,
		IdempotencyKey: uuid.MustParse("22222222-2222-2222-2222-222222222222"),
		Status:         domain.MintRequestStatusFailed,
		RetryCount:     RetryCount,
		MaxRetries:     5,
	}
}

func SuccessfulMintResult() domain.MintResult {
	return domain.MintResult{
		TransactionHash: "0xhash",
		BlockNumber:     12345,
		ErrorCode:       "",
	}
}

func RetryableMintResult() domain.MintResult {
	return domain.MintResult{
		ErrorCode:    "timeout",
		ErrorMessage: "rpc timeout",
		Retryable:    true,
	}
}

func PermanentFailMintResult() domain.MintResult {
	return domain.MintResult{
		ErrorCode: "invalid_address",
		Retryable: false,
	}
}
