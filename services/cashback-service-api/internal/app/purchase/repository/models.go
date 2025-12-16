package repository

import (
	"time"

	"github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/domain"
	"github.com/google/uuid"
)

// purchaseModel represents the database model for purchases
type purchaseModel struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID     uuid.UUID `gorm:"type:uuid;not null;index"`
	Amount     float64   `gorm:"not null"`
	MerchantID string    `gorm:"not null"`
	Status     string    `gorm:"not null;default:'pending'"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}

func (purchaseModel) TableName() string {
	return "purchases"
}

// toDomain converts database model to domain entity
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

// fromDomain converts domain entity to database model
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
