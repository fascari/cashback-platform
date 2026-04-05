package transaction

import (
	"time"

	"github.com/google/uuid"

	"github.com/cashback-platform/services/blockchain-adapter/internal/domain"
)

type transactionModel struct {
	ID              int64     `gorm:"primaryKey;autoIncrement"`
	IdempotencyKey  uuid.UUID `gorm:"type:uuid;uniqueIndex;not null"`
	WalletAddress   string    `gorm:"type:varchar(42);not null"`
	ChainID         string    `gorm:"type:varchar(50);not null;default:'ethereum'"`
	TokenAmount     string    `gorm:"type:varchar(78);not null"`
	TransactionHash string    `gorm:"type:varchar(66)"`
	BlockNumber     int64
	GasUsed         int64
	GasPrice        string                 `gorm:"type:varchar(78)"`
	Status          domain.TransactionStatus `gorm:"type:varchar(50);not null;default:'pending';index"`
	ErrorCode       string                 `gorm:"type:varchar(100)"`
	ErrorMessage    string                 `gorm:"type:text"`
	Nonce           int64
	CreatedAt       time.Time  `gorm:"autoCreateTime"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime"`
	ConfirmedAt     *time.Time
}

func (transactionModel) TableName() string {
	return "blockchain.blockchain_transactions"
}

func (m transactionModel) toDomain() domain.BlockchainTransaction {
	return domain.BlockchainTransaction{
		ID:              m.ID,
		IdempotencyKey:  m.IdempotencyKey,
		WalletAddress:   m.WalletAddress,
		ChainID:         m.ChainID,
		TokenAmount:     m.TokenAmount,
		TransactionHash: m.TransactionHash,
		BlockNumber:     m.BlockNumber,
		GasUsed:         m.GasUsed,
		GasPrice:        m.GasPrice,
		Status:          m.Status,
		ErrorCode:       m.ErrorCode,
		ErrorMessage:    m.ErrorMessage,
		Nonce:           m.Nonce,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
		ConfirmedAt:     m.ConfirmedAt,
	}
}

func fromDomain(tx domain.BlockchainTransaction) transactionModel {
	return transactionModel{
		ID:              tx.ID,
		IdempotencyKey:  tx.IdempotencyKey,
		WalletAddress:   tx.WalletAddress,
		ChainID:         tx.ChainID,
		TokenAmount:     tx.TokenAmount,
		TransactionHash: tx.TransactionHash,
		BlockNumber:     tx.BlockNumber,
		GasUsed:         tx.GasUsed,
		GasPrice:        tx.GasPrice,
		Status:          tx.Status,
		ErrorCode:       tx.ErrorCode,
		ErrorMessage:    tx.ErrorMessage,
		Nonce:           tx.Nonce,
		ConfirmedAt:     tx.ConfirmedAt,
	}
}


