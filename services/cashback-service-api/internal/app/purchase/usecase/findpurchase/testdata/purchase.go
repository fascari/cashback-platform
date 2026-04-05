package testdata

import purchasedomain "github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/domain"

const PurchaseID int64 = 42

func FoundPurchase() purchasedomain.Purchase {
	return purchasedomain.Purchase{ID: PurchaseID, UserID: 1, Amount: 200.0}
}
