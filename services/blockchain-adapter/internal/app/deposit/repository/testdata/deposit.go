//go:build integration

package testdata

import (
	"time"

	"github.com/cashback-platform/services/blockchain-adapter/internal/app/deposit/domain"
)

const (
	DepositID   int64 = 1
	ChainID           = "ethereum"
	TxHash            = "0xabc000"
	NewTxHash         = "0xdef111"
	BlockNumber int64 = 100
)

var FixtureTime = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

func NewDeposit() domain.Deposit {
	return domain.Deposit{
		ChainID:         ChainID,
		TransactionHash: NewTxHash,
		WalletAddress:   "0x9999",
		TokenAmount:     "500000000000000000",
		BlockNumber:     200,
		Status:          domain.StatusPending,
		DetectedAt:      FixtureTime,
	}
}

func ExistingDeposit() domain.Deposit {
	return domain.Deposit{
		ID:              DepositID,
		ChainID:         ChainID,
		TransactionHash: TxHash,
		WalletAddress:   "0x1234",
		TokenAmount:     "1000000000000000000",
		BlockNumber:     BlockNumber,
		Status:          domain.StatusPending,
		DetectedAt:      FixtureTime,
	}
}
