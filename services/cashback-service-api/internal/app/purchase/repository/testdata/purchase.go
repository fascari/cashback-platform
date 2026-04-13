//go:build integration

package testdata

import (
	"time"

	purchasedomain "github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/domain"
)

const (
	UserID            int64 = 1
	NewUserID         int64 = 2
	PurchaseID        int64 = 1
	AnotherPurchaseID int64 = 2
	NewPurchaseID     int64 = 10001
)

var (
	FixtureTime = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	FixedTime   = time.Date(2024, 6, 15, 10, 0, 0, 0, time.UTC)
)

func NewPurchase(userID int64) purchasedomain.Purchase {
	return purchasedomain.Purchase{
		UserID:     userID,
		Amount:     150.00,
		MerchantID: "merchant-abc",
		Status:     purchasedomain.StatusPending,
	}
}

func AnotherPurchase(userID int64) purchasedomain.Purchase {
	return purchasedomain.Purchase{
		UserID:     userID,
		Amount:     75.50,
		MerchantID: "merchant-xyz",
		Status:     purchasedomain.StatusPending,
	}
}

func CreatedPurchase() purchasedomain.Purchase {
	return purchasedomain.Purchase{
		ID:         NewPurchaseID,
		UserID:     UserID,
		Amount:     150.00,
		MerchantID: "merchant-abc",
		Status:     purchasedomain.StatusPending,
		CreatedAt:  FixedTime,
		UpdatedAt:  FixedTime,
	}
}

func ExistingPurchase() purchasedomain.Purchase {
	return purchasedomain.Purchase{
		ID:         PurchaseID,
		UserID:     UserID,
		Amount:     150.00,
		MerchantID: "merchant-abc",
		Status:     purchasedomain.StatusPending,
		CreatedAt:  FixtureTime,
		UpdatedAt:  FixtureTime,
	}
}
