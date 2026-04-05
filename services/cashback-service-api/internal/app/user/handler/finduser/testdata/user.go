package testdata

import "github.com/cashback-platform/services/cashback-service-api/internal/app/user/domain"

const ValidID = int64(1)

func ExistingUser() domain.User {
	return domain.User{ID: 1, ExternalID: "ext-123", Email: "user@example.com", WalletAddress: "0x1234567890abcdef1234"}
}
