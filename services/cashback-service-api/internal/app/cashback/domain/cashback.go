package domain

import (
	"time"

	"github.com/cashback-platform/kit/apperror"
)

const (
	ErrCodeInvalidUserID         = "error_invalid_cashback_user_id"
	ErrCodeInvalidPurchaseID     = "error_invalid_cashback_purchase_id"
	ErrCodeInvalidAmount         = "error_invalid_cashback_amount"
	ErrCodeInvalidPercentage     = "error_invalid_cashback_percentage"
	ErrCodeCashbackNotFound      = "error_cashback_not_found"
	ErrCodeCashbackAlreadyExists = "error_cashback_already_exists"

	StatusPending  Status = "pending"
	StatusApproved Status = "approved"
	StatusMinted   Status = "minted"
	StatusFailed   Status = "failed"
)

var (
	ErrInvalidUserID         = apperror.New(ErrCodeInvalidUserID, "invalid user ID")
	ErrInvalidPurchaseID     = apperror.New(ErrCodeInvalidPurchaseID, "invalid purchase ID")
	ErrInvalidAmount         = apperror.New(ErrCodeInvalidAmount, "invalid cashback amount")
	ErrInvalidPercentage     = apperror.New(ErrCodeInvalidPercentage, "invalid cashback percentage")
	ErrCashbackNotFound      = apperror.New(ErrCodeCashbackNotFound, "cashback not found")
	ErrCashbackAlreadyExists = apperror.New(ErrCodeCashbackAlreadyExists, "cashback already exists for this purchase")
)

type (
	Status string

	Cashback struct {
		ID              int64
		UserID          int64
		PurchaseID      int64
		Amount          float64
		CashbackPercent float64
		Status          Status
		CreatedAt       time.Time
		UpdatedAt       time.Time
	}
)

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
	return Cashback{
		UserID:          userID,
		PurchaseID:      purchaseID,
		Amount:          purchaseAmount * (cashbackPercent / 100),
		CashbackPercent: cashbackPercent,
		Status:          StatusPending,
		CreatedAt:       now,
		UpdatedAt:       now,
	}, nil
}

func (c Cashback) Approve() Cashback {
	c.Status = StatusApproved
	c.UpdatedAt = time.Now().UTC()
	return c
}

func (c Cashback) MarkAsFailed() Cashback {
	c.Status = StatusFailed
	c.UpdatedAt = time.Now().UTC()
	return c
}
