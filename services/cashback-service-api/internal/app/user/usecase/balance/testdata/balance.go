package testdata

import (
	userdomain "github.com/cashback-platform/services/cashback-service-api/internal/app/user/domain"
	balanceuc "github.com/cashback-platform/services/cashback-service-api/internal/app/user/usecase/balance"
)

const (
	UserID        int64 = 1
	WalletAddress       = "0xABCDEF1234567890"
)

func FoundUser() userdomain.User {
	return userdomain.User{ID: UserID, WalletAddress: WalletAddress, Email: "user@example.com"}
}

func TokenBalance() balanceuc.TokenBalance {
	return balanceuc.TokenBalance{
		WalletAddress: WalletAddress,
		Amount:        "1000000000000000000",
		BlockNumber:   100,
	}
}

func ZeroTokenBalance() balanceuc.TokenBalance {
	return balanceuc.TokenBalance{
		WalletAddress: WalletAddress,
		Amount:        "0",
		BlockNumber:   100,
	}
}

func InvalidWeiTokenBalance() balanceuc.TokenBalance {
	return balanceuc.TokenBalance{
		WalletAddress: WalletAddress,
		Amount:        "not-a-number",
		BlockNumber:   100,
	}
}
