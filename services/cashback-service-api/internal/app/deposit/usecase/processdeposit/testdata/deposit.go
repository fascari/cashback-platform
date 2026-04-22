package testdata

import (
	"time"

	"github.com/cashback-platform/services/cashback-service-api/internal/app/deposit/domain"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/deposit/usecase/processdeposit"
)

const (
	TxHash      = "0xabc123"
	FromAddress = "0xABCDEF1234567890ABCD"
	TokenAmount = "1000000000000000000" // 1 token in wei
	ChainID     = "ethereum-sepolia"
)

// ValidInput returns a well-formed Input for a new deposit event.
func ValidInput() processdeposit.Input {
	return processdeposit.Input{
		ChainID:         ChainID,
		TransactionHash: TxHash,
		FromAddress:     FromAddress,
		ToAddress:       "0xPlatformWallet",
		TokenAmount:     TokenAmount,
		BlockNumber:     12345678,
		DetectedAt:      time.Date(2024, 1, 15, 10, 45, 0, 0, time.UTC),
	}
}

// SavedReceipt returns a DepositReceipt as it would be returned after persistence.
func SavedReceipt() domain.DepositReceipt {
	return domain.DepositReceipt{
		ID:          1,
		UserID:      UserID,
		TxHash:      TxHash,
		FromAddress: FromAddress,
		Amount:      TokenAmount,
		ChainID:     ChainID,
		BlockNumber: 12345678,
		DetectedAt:  time.Date(2024, 1, 15, 10, 45, 0, 0, time.UTC),
	}
}
