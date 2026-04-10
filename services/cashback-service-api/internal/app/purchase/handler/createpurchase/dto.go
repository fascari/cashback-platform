package createpurchase

import (
	"strconv"

	"github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/domain"
	"github.com/cashback-platform/kit/validator"
)

type (
	InputPayload struct {
		UserID   string  `json:"user_id"  validate:"required,numeric"`
		Amount   float64 `json:"amount"   validate:"gt=0"`
		Merchant string  `json:"merchant" validate:"required"`
	}

	OutputPayload struct {
		ID         string  `json:"id"`
		UserID     string  `json:"user_id"`
		Amount     float64 `json:"amount"`
		MerchantID string  `json:"merchant_id"`
		Status     string  `json:"status"`
		CreatedAt  string  `json:"created_at"`
	}
)

func (p InputPayload) Validate() error {
	return validator.Validate(p)
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
