package testdata

import userdomain "github.com/cashback-platform/services/cashback-service-api/internal/app/user/domain"

const UserID int64 = 10

// FoundUser returns a User with a wallet address matching the test deposit.
func FoundUser() userdomain.User {
	return userdomain.User{
		ID:            UserID,
		WalletAddress: FromAddress,
	}
}
