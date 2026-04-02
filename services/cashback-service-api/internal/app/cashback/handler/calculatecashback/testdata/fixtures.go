package testdata

import (
	"github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/domain"
	calculatecashbackuc "github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/usecase/calculatecashback"
)

func ExistingCashback() domain.Cashback {
	return domain.Cashback{ID: 1, UserID: 42, PurchaseID: 1, Amount: 5.0, CashbackPercent: 5.0, Status: domain.StatusApproved}
}

func ValidPurchase() calculatecashbackuc.Purchase {
	return calculatecashbackuc.Purchase{ID: 1, UserID: 42, Amount: 100.0}
}

func ValidUser() calculatecashbackuc.User {
	return calculatecashbackuc.User{WalletAddress: "0xabc"}
}
