//go:build integration

package testdata

import (
	"time"

	cashdomain "github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/domain"
)

var (
	FixtureTime = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	FixedTime   = time.Date(2024, 6, 15, 10, 0, 0, 0, time.UTC)
)

const (
	UserID            int64 = 1
	PurchaseID        int64 = 1
	AnotherPurchaseID int64 = 2
	NewPurchaseID     int64 = 3
	CashbackID        int64 = 1
	MintedCashbackID  int64 = 2
	NewCashbackID     int64 = 10001
)

func NewCashback() cashdomain.Cashback {
	return cashdomain.Cashback{
		UserID:          UserID,
		PurchaseID:      NewPurchaseID,
		Amount:          15.0,
		CashbackPercent: 10.0,
		Status:          cashdomain.StatusPending,
	}
}

func CreatedCashback() cashdomain.Cashback {
	return cashdomain.Cashback{
		ID:              NewCashbackID,
		UserID:          UserID,
		PurchaseID:      NewPurchaseID,
		Amount:          15.0,
		CashbackPercent: 10.0,
		Status:          cashdomain.StatusPending,
		CreatedAt:       FixedTime,
		UpdatedAt:       FixedTime,
	}
}

func PendingCashback() cashdomain.Cashback {
	return cashdomain.Cashback{
		ID:              CashbackID,
		UserID:          UserID,
		PurchaseID:      PurchaseID,
		Amount:          15.0,
		CashbackPercent: 10.0,
		Status:          cashdomain.StatusPending,
		CreatedAt:       FixtureTime,
		UpdatedAt:       FixtureTime,
	}
}

func MintedCashback() cashdomain.Cashback {
	return cashdomain.Cashback{
		ID:              MintedCashbackID,
		UserID:          UserID,
		PurchaseID:      AnotherPurchaseID,
		Amount:          10.0,
		CashbackPercent: 10.0,
		Status:          cashdomain.StatusMinted,
		CreatedAt:       FixtureTime,
		UpdatedAt:       FixtureTime,
	}
}
