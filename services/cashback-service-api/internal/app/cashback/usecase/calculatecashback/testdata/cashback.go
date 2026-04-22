package testdata

import cashdomain "github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/domain"

const (
	PurchaseID int64 = 1
	UserID     int64 = 10
)

func pInt64(v int64) *int64 { return &v }

func ApprovedCashback() cashdomain.Cashback {
	return cashdomain.Cashback{ID: 1, UserID: UserID, PurchaseID: pInt64(PurchaseID), Amount: 5.0, Status: cashdomain.StatusApproved}
}

func ExistingCashback() cashdomain.Cashback {
	return cashdomain.Cashback{ID: 99, PurchaseID: pInt64(PurchaseID)}
}
