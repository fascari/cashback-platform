package findusercashback

import (
	"github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/domain"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/usecase/findusercashback"
)

type (
	CashbackItem struct {
		ID              string  `json:"id"`
		PurchaseID      string  `json:"purchase_id"`
		Amount          float64 `json:"amount"`
		CashbackPercent float64 `json:"cashback_percent"`
		Status          string  `json:"status"`
		CreatedAt       string  `json:"created_at"`
	}

	OutputPayload struct {
		UserID         string         `json:"user_id"`
		Cashbacks      []CashbackItem `json:"cashbacks"`
		TotalMinted    float64        `json:"total_minted"`
		TotalCashbacks int            `json:"total_cashbacks"`
	}
)

func ToOutputPayload(summary findusercashback.UserCashbackSummary) OutputPayload {
	cashbacks := make([]CashbackItem, len(summary.Cashbacks))
	for i, c := range summary.Cashbacks {
		cashbacks[i] = toCashbackItem(c)
	}

	return OutputPayload{
		UserID:         summary.UserID.String(),
		Cashbacks:      cashbacks,
		TotalMinted:    summary.TotalMinted,
		TotalCashbacks: summary.TotalCashbacks,
	}
}

func toCashbackItem(c domain.Cashback) CashbackItem {
	return CashbackItem{
		ID:              c.ID.String(),
		PurchaseID:      c.PurchaseID.String(),
		Amount:          c.Amount,
		CashbackPercent: c.CashbackPercent,
		Status:          c.Status,
		CreatedAt:       c.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
