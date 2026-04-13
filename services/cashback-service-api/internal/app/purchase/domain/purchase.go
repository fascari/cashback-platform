package domain

import (
	"time"

	"github.com/cashback-platform/kit/apperror"
	"github.com/cashback-platform/kit/clock"
)

const (
	ErrCodeInvalidAmount    = "error_invalid_purchase_amount"
	ErrCodeInvalidUserID    = "error_invalid_purchase_user_id"
	ErrCodeInvalidMerchant  = "error_invalid_purchase_merchant"
	ErrCodePurchaseNotFound = "error_purchase_not_found"

	StatusPending Status = "pending"
)

var (
	ErrInvalidAmount    = apperror.New(ErrCodeInvalidAmount, "invalid purchase amount")
	ErrInvalidUserID    = apperror.New(ErrCodeInvalidUserID, "invalid user ID")
	ErrInvalidMerchant  = apperror.New(ErrCodeInvalidMerchant, "invalid merchant ID")
	ErrPurchaseNotFound = apperror.New(ErrCodePurchaseNotFound, "purchase not found")
)

type (
	Status string

	Purchase struct {
		ID         int64
		UserID     int64
		Amount     float64
		MerchantID string
		Status     Status
		CreatedAt  time.Time
		UpdatedAt  time.Time
	}
)

func NewPurchase(userID int64, amount float64, merchant string) Purchase {
	now := clock.Now().UTC()
	return Purchase{
		UserID:     userID,
		Amount:     amount,
		MerchantID: merchant,
		Status:     StatusPending,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}
