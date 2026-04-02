package repository

import (
	"time"

	"github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/domain"
)

type cashbackModel struct {
	ID              int64   `gorm:"primaryKey;autoIncrement"`
	UserID          int64   `gorm:"not null;index"`
	PurchaseID      int64   `gorm:"not null;uniqueIndex"`
	Amount          float64 `gorm:"not null"`
	CashbackPercent float64 `gorm:"not null"`
	Status          string  `gorm:"not null;default:'pending';index"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (cashbackModel) TableName() string { return "cashback_ledger" }

func (m cashbackModel) toDomain() domain.Cashback {
	return domain.Cashback{
		ID:              m.ID,
		UserID:          m.UserID,
		PurchaseID:      m.PurchaseID,
		Amount:          m.Amount,
		CashbackPercent: m.CashbackPercent,
		Status:          m.Status,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

func fromDomain(cashback domain.Cashback) cashbackModel {
	return cashbackModel{
		ID:              cashback.ID,
		UserID:          cashback.UserID,
		PurchaseID:      cashback.PurchaseID,
		Amount:          cashback.Amount,
		CashbackPercent: cashback.CashbackPercent,
		Status:          cashback.Status,
		CreatedAt:       cashback.CreatedAt,
		UpdatedAt:       cashback.UpdatedAt,
	}
}
