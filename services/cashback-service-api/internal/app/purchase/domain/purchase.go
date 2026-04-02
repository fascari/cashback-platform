// Package domain contains the core business entities and rules for purchases.
package domain

import (
	"errors"
	"time"
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
	ID         int64
	UserID     int64
	Amount     float64
	MerchantID string
	Status     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// NewPurchase creates a new purchase instance.
// Status is initialized as "pending" by default. The ID is set by the database on insert.
func NewPurchase(userID int64, amount float64, merchant string) Purchase {
	now := time.Now().UTC()
	return Purchase{
		UserID:     userID,
		Amount:     amount,
		MerchantID: merchant,
		Status:     "pending",
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}
