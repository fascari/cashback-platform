package repository

import (
	"time"

	"github.com/cashback-platform/services/blockchain-adapter/internal/app/deposit/domain"
)

type depositModel struct {
	ID              int64      `gorm:"primaryKey;autoIncrement"`
	ChainID         string     `gorm:"column:chain_id;not null"`
	TransactionHash string     `gorm:"column:transaction_hash;not null"`
	WalletAddress   string     `gorm:"column:wallet_address;not null"`
	TokenAmount     string     `gorm:"column:token_amount;not null"`
	BlockNumber     int64      `gorm:"column:block_number;not null"`
	Status          string     `gorm:"column:status;not null;default:pending"`
	DetectedAt      time.Time  `gorm:"column:detected_at"`
	ProcessedAt     *time.Time `gorm:"column:processed_at"`
}

func (depositModel) TableName() string {
	return "blockchain.detected_deposits"
}

func (m depositModel) ToDomain() domain.Deposit {
	return domain.Deposit{
		ID:              m.ID,
		ChainID:         m.ChainID,
		TransactionHash: m.TransactionHash,
		WalletAddress:   m.WalletAddress,
		TokenAmount:     m.TokenAmount,
		BlockNumber:     m.BlockNumber,
		Status:          domain.Status(m.Status),
		DetectedAt:      m.DetectedAt,
		ProcessedAt:     m.ProcessedAt,
	}
}

func fromDomain(d domain.Deposit) depositModel {
	return depositModel{
		ChainID:         d.ChainID,
		TransactionHash: d.TransactionHash,
		WalletAddress:   d.WalletAddress,
		TokenAmount:     d.TokenAmount,
		BlockNumber:     d.BlockNumber,
		Status:          string(d.Status),
		DetectedAt:      d.DetectedAt,
	}
}
