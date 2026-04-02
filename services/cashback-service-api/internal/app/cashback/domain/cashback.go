// Package domain contains the core business entities and rules for cashback.
package domain

import (
	"errors"
	"time"
)

// Cashback status values represent the lifecycle of a cashback transaction.
const (
	StatusPending  = "pending"
	StatusApproved = "approved"
	StatusMinted   = "minted"
	StatusFailed   = "failed"
)

// Sentinel errors for cashback domain validation.
var (
	ErrInvalidUserID     = errors.New("invalid user ID")
	ErrInvalidPurchaseID = errors.New("invalid purchase ID")
	ErrInvalidAmount     = errors.New("invalid cashback amount")
	ErrInvalidPercentage = errors.New("invalid cashback percentage")
	ErrCashbackNotFound  = errors.New("cashback not found")
)

// Cashback represents a cashback transaction in the system.
// It tracks the cashback amount, status, and relationships to users and purchases.
type Cashback struct {
	ID              int64
	UserID          int64
	PurchaseID      int64
	Amount          float64
	CashbackPercent float64
	Status          string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// NewCashback creates a new cashback instance with validation.
// It calculates the cashback amount based on the purchase amount and percentage.
// Returns an error if any validation fails. The ID is set by the database on insert.
func NewCashback(userID, purchaseID int64, purchaseAmount, cashbackPercent float64) (Cashback, error) {
	if userID == 0 {
		return Cashback{}, ErrInvalidUserID
	}
	if purchaseID == 0 {
		return Cashback{}, ErrInvalidPurchaseID
	}
	if purchaseAmount <= 0 {
		return Cashback{}, ErrInvalidAmount
	}
	if cashbackPercent <= 0 || cashbackPercent > 100 {
		return Cashback{}, ErrInvalidPercentage
	}

	now := time.Now().UTC()
	cashbackAmount := purchaseAmount * (cashbackPercent / 100)

	return Cashback{
		UserID:          userID,
		PurchaseID:      purchaseID,
		Amount:          cashbackAmount,
		CashbackPercent: cashbackPercent,
		Status:          StatusPending,
		CreatedAt:       now,
		UpdatedAt:       now,
	}, nil
}

// Approve transitions the cashback to approved status.
// This indicates the cashback is ready to be minted as tokens.
func (c *Cashback) Approve() {
	c.Status = StatusApproved
	c.UpdatedAt = time.Now().UTC()
}

// MarkAsMinted transitions the cashback to minted status.
// This indicates tokens have been successfully minted on the blockchain.
func (c *Cashback) MarkAsMinted() {
	c.Status = StatusMinted
	c.UpdatedAt = time.Now().UTC()
}

// MarkAsFailed transitions the cashback to failed status.
// This indicates the minting process failed and may require manual intervention.
func (c *Cashback) MarkAsFailed() {
	c.Status = StatusFailed
	c.UpdatedAt = time.Now().UTC()
}
