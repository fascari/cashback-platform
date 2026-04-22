package domain

import (
	"time"

	"github.com/cashback-platform/kit/apperror"
	"github.com/cashback-platform/kit/clock"
)

const (
	ErrCodeDepositNotFound         = "error_deposit_not_found"
	ErrCodeDepositAlreadyProcessed = "error_deposit_already_processed"
	ErrCodeDepositInvalidUser      = "error_deposit_invalid_user"
	ErrCodeDepositInvalidTxHash    = "error_deposit_invalid_tx_hash"
)

var (
	ErrDepositNotFound         = apperror.New(ErrCodeDepositNotFound, "deposit not found")
	ErrDepositAlreadyProcessed = apperror.New(ErrCodeDepositAlreadyProcessed, "deposit already processed")
	ErrDepositInvalidUser      = apperror.New(ErrCodeDepositInvalidUser, "invalid user for deposit")
	ErrDepositInvalidTxHash    = apperror.New(ErrCodeDepositInvalidTxHash, "invalid transaction hash")
)

type DepositReceipt struct {
	ID          int64
	UserID      int64
	TxHash      string
	FromAddress string
	Amount      string
	ChainID     string
	BlockNumber int64
	DetectedAt  time.Time
	CreatedAt   time.Time
}

func NewDepositReceipt(userID int64, txHash, fromAddress, amount, chainID string, blockNumber int64, detectedAt time.Time) (DepositReceipt, error) {
	if userID == 0 {
		return DepositReceipt{}, ErrDepositInvalidUser
	}
	if txHash == "" {
		return DepositReceipt{}, ErrDepositInvalidTxHash
	}
	return DepositReceipt{
		UserID:      userID,
		TxHash:      txHash,
		FromAddress: fromAddress,
		Amount:      amount,
		ChainID:     chainID,
		BlockNumber: blockNumber,
		DetectedAt:  detectedAt.UTC(),
		CreatedAt:   clock.Now().UTC(),
	}, nil
}
