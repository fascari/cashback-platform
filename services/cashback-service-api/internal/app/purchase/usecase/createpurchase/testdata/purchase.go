package testdata

import purchasedomain "github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/domain"

const (
	UserID     int64   = 1
	Amount     float64 = 150.0
	MerchantID         = "store-123"
)

func CreatedPurchase() purchasedomain.Purchase {
	return purchasedomain.Purchase{ID: 99, UserID: UserID, Amount: Amount, MerchantID: MerchantID}
}
