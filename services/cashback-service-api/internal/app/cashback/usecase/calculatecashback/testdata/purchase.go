package testdata

import "github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/usecase/calculatecashback"

func NewPurchase() calculatecashback.Purchase {
	return calculatecashback.Purchase{ID: PurchaseID, UserID: UserID, Amount: 100.0}
}
