// Package domain contains the core business entities and rules for purchases.
package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Sentinel errors for purchase domain validation.
var (
	ErrInvalidAmount   = errors.New("invalid purchase amount")
	ErrInvalidUserID   = errors.New("invalid user ID")
	ErrInvalidMerchant = errors.New("invalid merchant ID")
)

// Purchase represents a purchase transaction in the system.
// It tracks the purchase amount, merchant, and user relationship.
type Purchase struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	Amount     float64
	MerchantID string
	Status     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// NewPurchase creates a new purchase instance.
// Status is initialized as "pending" by default.
func NewPurchase(userID uuid.UUID, amount float64, merchant string) Purchase {
	now := time.Now().UTC()
	return Purchase{
		ID:         uuid.New(),
		UserID:     userID,
		Amount:     amount,
		MerchantID: merchant,
		Status:     "pending",
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}
