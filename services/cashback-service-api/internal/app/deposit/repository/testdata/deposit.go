//go:build integration

package testdata

import (
	"time"

	"github.com/cashback-platform/services/cashback-service-api/internal/app/deposit/domain"
)

const (
	DepositReceiptID = int64(1)
	TxHash           = "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890ab"
	FromAddress      = "0xABCDEF1234567890ABCD"
	UserID           = int64(1)
)

// NewDepositReceipt returns an unsaved DepositReceipt ready for insertion.
func NewDepositReceipt() domain.DepositReceipt {
	return domain.DepositReceipt{
		UserID:      UserID,
		TxHash:      TxHash,
		FromAddress: FromAddress,
		Amount:      "1000000000000000000",
		ChainID:     "ethereum-sepolia",
		BlockNumber: 12345678,
		DetectedAt:  time.Date(2024, 1, 15, 10, 45, 0, 0, time.UTC),
	}
}

// SavedDepositReceipt returns the persisted form of NewDepositReceipt as it appears in fixtures.
func SavedDepositReceipt() domain.DepositReceipt {
	r := NewDepositReceipt()
	r.ID = DepositReceiptID
	r.CreatedAt = time.Date(2024, 1, 15, 10, 45, 0, 0, time.UTC)
	return r
}
