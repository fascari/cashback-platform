package testdata

import (
	"github.com/google/uuid"

	"github.com/cashback-platform/services/mint-consumer/internal/app/mint/domain"
	"github.com/cashback-platform/services/mint-consumer/internal/app/mint/usecase/mintcashback"
)

const (
	CashbackID    int64 = 1
	UserID        int64 = 10
	WalletAddress       = "0xABCDEF1234567890"
	TokenAmount         = "100"
	ChainID             = "sepolia"
)

var EventID = uuid.MustParse("11111111-1111-1111-1111-111111111111")

func NewInput() mintcashback.Input {
	return mintcashback.Input{
		EventID:       EventID,
		CashbackID:    CashbackID,
		UserID:        UserID,
		WalletAddress: WalletAddress,
		TokenAmount:   TokenAmount,
		ChainID:       ChainID,
	}
}

func PendingMintRequest() domain.MintRequest {
	return domain.MintRequest{
		ID:            42,
		CashbackID:    CashbackID,
		UserID:        UserID,
		WalletAddress: WalletAddress,
		TokenAmount:   TokenAmount,
		Status:        domain.MintRequestStatusPending,
		MaxRetries:    5,
	}
}

func CompletedMintRequest() domain.MintRequest {
	return domain.MintRequest{
		ID:              42,
		CashbackID:      CashbackID,
		UserID:          UserID,
		WalletAddress:   WalletAddress,
		TokenAmount:     TokenAmount,
		Status:          domain.MintRequestStatusCompleted,
		TransactionHash: "0xhash",
		BlockNumber:     12345,
		MaxRetries:      5,
	}
}

func SuccessfulMintResult() domain.MintResult {
	return domain.MintResult{
		TransactionHash: "0xhash",
		BlockNumber:     12345,
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
		ErrorCode:    "invalid_address",
		ErrorMessage: "bad address",
		Retryable:    false,
	}
}
