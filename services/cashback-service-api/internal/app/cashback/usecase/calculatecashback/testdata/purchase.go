package testdata

import purchasedomain "github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/domain"

func NewPurchase() purchasedomain.Purchase {
	return purchasedomain.Purchase{ID: PurchaseID, UserID: UserID, Amount: 100.0}
}
