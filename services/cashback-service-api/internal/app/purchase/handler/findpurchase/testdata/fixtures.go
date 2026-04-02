package testdata

import "github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/domain"

func ExistingPurchase() domain.Purchase {
	return domain.Purchase{ID: 1, UserID: 42, Amount: 100.0, MerchantID: "shopA", Status: domain.StatusPending}
}
