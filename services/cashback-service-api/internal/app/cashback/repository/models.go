package repository

import (
	"time"

	"github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/domain"
)

type cashbackModel struct {
	ID               int64         `gorm:"primaryKey;autoIncrement"`
	UserID           int64         `gorm:"not null;index"`
	PurchaseID       *int64        `gorm:"index"`
	DepositReceiptID *int64        `gorm:"index"`
	Amount           float64       `gorm:"not null"`
	CashbackPercent  float64       `gorm:"not null"`
	Status           domain.Status `gorm:"not null;default:'pending';index"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (cashbackModel) TableName() string {
	return "cashback.cashback_ledger"
}

func (m cashbackModel) toDomain() domain.Cashback {
	return domain.Cashback{
		ID:               m.ID,
		UserID:           m.UserID,
		PurchaseID:       m.PurchaseID,
		DepositReceiptID: m.DepositReceiptID,
		Amount:           m.Amount,
		CashbackPercent:  m.CashbackPercent,
		Status:           m.Status,
		CreatedAt:        m.CreatedAt.UTC(),
		UpdatedAt:        m.UpdatedAt.UTC(),
	}
}

func fromDomain(c domain.Cashback) cashbackModel {
	return cashbackModel{
		ID:               c.ID,
		UserID:           c.UserID,
		PurchaseID:       c.PurchaseID,
		DepositReceiptID: c.DepositReceiptID,
		Amount:           c.Amount,
		CashbackPercent:  c.CashbackPercent,
		Status:           c.Status,
		CreatedAt:        c.CreatedAt,
		UpdatedAt:        c.UpdatedAt,
	}
}
