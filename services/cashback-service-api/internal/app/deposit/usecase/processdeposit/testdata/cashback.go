package testdata

import cashdomain "github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/domain"

func pInt64(v int64) *int64 { return &v }

// ApprovedCashback returns a Cashback as it would be returned after creation and approval.
func ApprovedCashback() cashdomain.Cashback {
	return cashdomain.Cashback{
		ID:               1,
		UserID:           UserID,
		DepositReceiptID: pInt64(1),
		Amount:           0.01,
		CashbackPercent:  1.0,
		Status:           cashdomain.StatusApproved,
	}
}
