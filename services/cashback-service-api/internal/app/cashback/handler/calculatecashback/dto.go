package calculatecashback

import (
	"github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/domain"
)

type (
	InputPayload struct {
		PurchaseID string `json:"purchase_id"`
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
	if p.PurchaseID == "" {
		return domain.ErrInvalidPurchaseID
	}
	return nil
}

func ToOutputPayload(cashback domain.Cashback) OutputPayload {
	return OutputPayload{
		ID:              cashback.ID.String(),
		UserID:          cashback.UserID.String(),
		PurchaseID:      cashback.PurchaseID.String(),
		Amount:          cashback.Amount,
		CashbackPercent: cashback.CashbackPercent,
		Status:          cashback.Status,
		CreatedAt:       cashback.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
