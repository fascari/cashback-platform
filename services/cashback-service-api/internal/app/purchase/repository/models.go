package repository

import (
	"time"

	"github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/domain"
)

type purchaseModel struct {
	ID         int64         `gorm:"primaryKey;autoIncrement"`
	UserID     int64         `gorm:"not null;index"`
	Amount     float64       `gorm:"not null"`
	MerchantID string        `gorm:"not null"`
	Status     domain.Status `gorm:"not null;default:'pending'"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (purchaseModel) TableName() string {
	return "cashback.purchases"
}

func (m purchaseModel) toDomain() domain.Purchase {
	return domain.Purchase{
		ID:         m.ID,
		UserID:     m.UserID,
		Amount:     m.Amount,
		MerchantID: m.MerchantID,
		Status:     m.Status,
		CreatedAt:  m.CreatedAt.UTC(),
		UpdatedAt:  m.UpdatedAt.UTC(),
	}
}

func fromDomain(p domain.Purchase) purchaseModel {
	return purchaseModel{
		ID:         p.ID,
		UserID:     p.UserID,
		Amount:     p.Amount,
		MerchantID: p.MerchantID,
		Status:     p.Status,
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
	}
}
