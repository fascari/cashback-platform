package testdata

import "github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/domain"

func ExistingCashback() domain.Cashback {
	return domain.Cashback{ID: 1, UserID: 42, PurchaseID: 1, Amount: 5.0, CashbackPercent: 5.0, Status: domain.StatusApproved}
}
