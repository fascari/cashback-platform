package testdata

import (
	"github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/domain"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/usecase/findusercashback"
)

const ValidUserID = int64(42)

func pInt64(v int64) *int64 { return &v }

func CashbackSummary() findusercashback.UserCashbackSummary {
	return findusercashback.UserCashbackSummary{
		UserID:         42,
		TotalMinted:    10.0,
		TotalCashbacks: 2,
		Cashbacks: []domain.Cashback{
			{ID: 1, UserID: 42, PurchaseID: pInt64(1), Amount: 5.0, Status: domain.StatusMinted},
			{ID: 2, UserID: 42, PurchaseID: pInt64(2), Amount: 5.0, Status: domain.StatusMinted},
		},
	}
}
