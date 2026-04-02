package repository

import (
	"time"

	"github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/domain"
)

type purchaseModel struct {
	ID         int64   `gorm:"primaryKey;autoIncrement"`
	UserID     int64   `gorm:"not null;index"`
	Amount     float64 `gorm:"not null"`
	MerchantID string  `gorm:"not null"`
	Status     string  `gorm:"not null;default:'pending'"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (purchaseModel) TableName() string { return "purchases" }

func (m purchaseModel) toDomain() domain.Purchase {
	return domain.Purchase{
		ID:         m.ID,
		UserID:     m.UserID,
		Amount:     m.Amount,
		MerchantID: m.MerchantID,
		Status:     m.Status,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

func fromDomain(purchase domain.Purchase) purchaseModel {
	return purchaseModel{
		ID:         purchase.ID,
		UserID:     purchase.UserID,
		Amount:     purchase.Amount,
		MerchantID: purchase.MerchantID,
		Status:     purchase.Status,
		CreatedAt:  purchase.CreatedAt,
		UpdatedAt:  purchase.UpdatedAt,
	}
}
