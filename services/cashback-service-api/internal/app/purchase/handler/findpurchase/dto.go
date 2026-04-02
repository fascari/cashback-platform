package findpurchase

import (
	"strconv"

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
		ID:         strconv.FormatInt(purchase.ID, 10),
		UserID:     strconv.FormatInt(purchase.UserID, 10),
		Amount:     purchase.Amount,
		MerchantID: purchase.MerchantID,
		Status:     string(purchase.Status),
		CreatedAt:  purchase.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
