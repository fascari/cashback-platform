//go:build integration

package testdata

import (
	"time"

	userdomain "github.com/cashback-platform/services/cashback-service-api/internal/app/user/domain"
)

const (
	UserID    int64 = 1
	NewUserID int64 = 10001
)

var (
	FixtureTime = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	FixedTime   = time.Date(2024, 6, 15, 10, 0, 0, 0, time.UTC)
)

func NewUser() userdomain.User {
	return userdomain.User{
		ExternalID:    "ext-002",
		Email:         "other@example.com",
		WalletAddress: "0x1234567890ABCDEF1234",
	}
}

func AnotherUser() userdomain.User {
	return userdomain.User{
		ExternalID:    "ext-002",
		Email:         "other@example.com",
		WalletAddress: "0x1234567890ABCDEF1234",
	}
}

func CreatedUser() userdomain.User {
	return userdomain.User{
		ID:            NewUserID,
		ExternalID:    "ext-002",
		Email:         "other@example.com",
		WalletAddress: "0x1234567890ABCDEF1234",
		CreatedAt:     FixedTime,
		UpdatedAt:     FixedTime,
	}
}

func ExistingUser() userdomain.User {
	return userdomain.User{
		ID:            UserID,
		ExternalID:    "ext-001",
		Email:         "user@example.com",
		WalletAddress: "0xABCDEF1234567890ABCD",
		CreatedAt:     FixtureTime,
		UpdatedAt:     FixtureTime,
	}
}
