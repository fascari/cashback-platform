package repository

import (
	"time"

	"github.com/cashback-platform/services/cashback-service-api/internal/app/deposit/domain"
)

type depositReceiptModel struct {
	ID          int64     `gorm:"primaryKey;autoIncrement"`
	UserID      int64     `gorm:"not null;index:idx_deposit_receipts_user_id"`
	TxHash      string    `gorm:"not null;uniqueIndex:uq_deposit_receipts_tx_hash;size:66"`
	FromAddress string    `gorm:"not null;size:42"`
	Amount      string    `gorm:"not null;size:78"`
	ChainID     string    `gorm:"not null;size:32"`
	BlockNumber int64     `gorm:"not null"`
	DetectedAt  time.Time `gorm:"not null"`
	CreatedAt   time.Time `gorm:"not null;autoCreateTime"`
}

func (depositReceiptModel) TableName() string {
	return "cashback.deposit_receipts"
}

func (m depositReceiptModel) toDomain() domain.DepositReceipt {
	return domain.DepositReceipt{
		ID:          m.ID,
		UserID:      m.UserID,
		TxHash:      m.TxHash,
		FromAddress: m.FromAddress,
		Amount:      m.Amount,
		ChainID:     m.ChainID,
		BlockNumber: m.BlockNumber,
		DetectedAt:  m.DetectedAt.UTC(),
		CreatedAt:   m.CreatedAt.UTC(),
	}
}

func fromDomain(r domain.DepositReceipt) depositReceiptModel {
	return depositReceiptModel{
		ID:          r.ID,
		UserID:      r.UserID,
		TxHash:      r.TxHash,
		FromAddress: r.FromAddress,
		Amount:      r.Amount,
		ChainID:     r.ChainID,
		BlockNumber: r.BlockNumber,
		DetectedAt:  r.DetectedAt,
		CreatedAt:   r.CreatedAt,
	}
}
