package testdata

import "github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/usecase/calculatecashback"

func NewUser() calculatecashback.User {
	return calculatecashback.User{WalletAddress: "0xABC"}
}
