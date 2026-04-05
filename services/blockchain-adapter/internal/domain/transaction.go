package domain

import (
	"time"

	"github.com/google/uuid"
)

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusSubmitted TransactionStatus = "submitted"
	TransactionStatusConfirmed TransactionStatus = "confirmed"
	TransactionStatusFailed    TransactionStatus = "failed"
)

type (
	TransactionStatus string

	BlockchainTransaction struct {
		ID              int64
		IdempotencyKey  uuid.UUID
		WalletAddress   string
		ChainID         string
		TokenAmount     string
		TransactionHash string
		BlockNumber     int64
		GasUsed         int64
		GasPrice        string
		Status          TransactionStatus
		ErrorCode       string
		ErrorMessage    string
		Nonce           int64
		CreatedAt       time.Time
		UpdatedAt       time.Time
		ConfirmedAt     *time.Time
	}
)

func (t BlockchainTransaction) IsFinalized() bool {
	return t.Status == TransactionStatusSubmitted || t.Status == TransactionStatusConfirmed
}

func (t BlockchainTransaction) IsFailed() bool {
	return t.Status == TransactionStatusFailed
}
