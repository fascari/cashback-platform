package testdata

import userdomain "github.com/cashback-platform/services/cashback-service-api/internal/app/user/domain"

func NewUser() userdomain.User {
	return userdomain.User{WalletAddress: "0xABC"}
}
