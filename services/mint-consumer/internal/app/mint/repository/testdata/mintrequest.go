//go:build integration

package testdata

import (
	"github.com/cashback-platform/services/mint-consumer/internal/app/mint/domain"
	"github.com/google/uuid"
)

const (
	CashbackID    int64 = 1
	UserID        int64 = 10
	WalletAddress       = "0xABCDEF1234567890ABCD"
	TokenAmount         = "100"
)

var (
	EventID        = uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	IdempotencyKey = uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
)

func NewPendingMintRequest() domain.MintRequest {
	return domain.MintRequest{
		CashbackID:     CashbackID,
		UserID:         UserID,
		WalletAddress:  WalletAddress,
		TokenAmount:    TokenAmount,
		IdempotencyKey: IdempotencyKey,
		Status:         domain.MintRequestStatusPending,
		MaxRetries:     5,
	}
}
