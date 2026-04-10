package calculatecashback

import (
	"strconv"

	"github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/domain"
	"github.com/cashback-platform/kit/validator"
)

type (
	InputPayload struct {
		PurchaseID string `json:"purchase_id" validate:"required,numeric"`
	}

	OutputPayload struct {
		ID              string  `json:"id"`
		UserID          string  `json:"user_id"`
		PurchaseID      string  `json:"purchase_id"`
		Amount          float64 `json:"amount"`
		CashbackPercent float64 `json:"cashback_percent"`
		Status          string  `json:"status"`
		CreatedAt       string  `json:"created_at"`
	}
)

func (p InputPayload) Validate() error {
	return validator.Validate(p)
}

func ToOutputPayload(cashback domain.Cashback) OutputPayload {
	return OutputPayload{
		ID:              strconv.FormatInt(cashback.ID, 10),
		UserID:          strconv.FormatInt(cashback.UserID, 10),
		PurchaseID:      strconv.FormatInt(cashback.PurchaseID, 10),
		Amount:          cashback.Amount,
		CashbackPercent: cashback.CashbackPercent,
		Status:          string(cashback.Status),
		CreatedAt:       cashback.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
