package testdata

import "github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/domain"

func pInt64(v int64) *int64 { return &v }

func ExistingCashback() domain.Cashback {
	return domain.Cashback{ID: 1, UserID: 42, PurchaseID: pInt64(1), Amount: 5.0, CashbackPercent: 5.0, Status: domain.StatusApproved}
}
