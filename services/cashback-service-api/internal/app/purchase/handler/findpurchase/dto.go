package findpurchase

import (
	"github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/domain"
)

type OutputPayload struct {
	ID         string  `json:"id"`
	UserID     string  `json:"user_id"`
	Amount     float64 `json:"amount"`
	MerchantID string  `json:"merchant_id"`
	Status     string  `json:"status"`
	CreatedAt  string  `json:"created_at"`
}

func ToOutputPayload(purchase domain.Purchase) OutputPayload {
	return OutputPayload{
		ID:         purchase.ID.String(),
		UserID:     purchase.UserID.String(),
		Amount:     purchase.Amount,
		MerchantID: purchase.MerchantID,
		Status:     purchase.Status,
		CreatedAt:  purchase.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
