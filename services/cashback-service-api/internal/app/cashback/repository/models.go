package repository

import (
	"time"

	"github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/domain"
	"github.com/google/uuid"
)

type cashbackModel struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID          uuid.UUID `gorm:"type:uuid;not null;index"`
	PurchaseID      uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	Amount          float64   `gorm:"not null"`
	CashbackPercent float64   `gorm:"not null"`
	Status          string    `gorm:"not null;default:'pending';index"`
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`
}

func (cashbackModel) TableName() string {
	return "cashback_ledger"
}

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
