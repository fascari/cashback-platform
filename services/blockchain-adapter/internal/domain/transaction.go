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
	// TransactionStatus represents the status of a blockchain transaction
	TransactionStatus string

	// BlockchainTransaction represents a blockchain transaction record
	BlockchainTransaction struct {
		ID              uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
		IdempotencyKey  uuid.UUID `gorm:"type:uuid;uniqueIndex;not null"`
		WalletAddress   string    `gorm:"type:varchar(42);not null"`
		TokenAmount     string    `gorm:"type:varchar(78);not null"`
		TransactionHash string    `gorm:"type:varchar(66)"`
		BlockNumber     int64
		GasUsed         int64
		GasPrice        string            `gorm:"type:varchar(78)"`
		Status          TransactionStatus `gorm:"type:varchar(50);not null;default:'pending';index"`
		ErrorCode       string            `gorm:"type:varchar(100)"`
		ErrorMessage    string            `gorm:"type:text"`
		Nonce           int64
		CreatedAt       time.Time `gorm:"autoCreateTime"`
		UpdatedAt       time.Time `gorm:"autoUpdateTime"`
		ConfirmedAt     *time.Time
	}
)

// TableName specifies the table name for GORM
func (BlockchainTransaction) TableName() string {
	return "blockchain_transactions"
}
