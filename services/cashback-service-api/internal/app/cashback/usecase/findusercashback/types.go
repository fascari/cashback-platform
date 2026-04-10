package findusercashback

import "github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/domain"

type UserCashbackSummary struct {
	UserID         int64
	Cashbacks      []domain.Cashback
	TotalMinted    float64
	TotalCashbacks int
}
