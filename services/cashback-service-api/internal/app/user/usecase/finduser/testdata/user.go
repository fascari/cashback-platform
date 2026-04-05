package testdata

import userdomain "github.com/cashback-platform/services/cashback-service-api/internal/app/user/domain"

const UserID int64 = 7

func FoundUser() userdomain.User {
	return userdomain.User{ID: UserID, Email: "user@example.com", WalletAddress: "0xABC"}
}
