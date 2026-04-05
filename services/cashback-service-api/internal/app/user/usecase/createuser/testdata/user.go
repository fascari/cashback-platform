package testdata

import userdomain "github.com/cashback-platform/services/cashback-service-api/internal/app/user/domain"

const (
	ExternalID    = "ext-1"
	Email         = "user@example.com"
	WalletAddress = "0xABC"
)

func CreatedUser() userdomain.User {
	return userdomain.User{ID: 1, ExternalID: ExternalID, Email: Email, WalletAddress: WalletAddress}
}
