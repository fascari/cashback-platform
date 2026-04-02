package testdata

import (
	userdomain "github.com/cashback-platform/services/cashback-service-api/internal/app/user/domain"
)

func CreatedUser() userdomain.User {
	p := ValidPayload()
	return userdomain.User{
		ID:            1,
		ExternalID:    p.ExternalID,
		Email:         p.Email,
		WalletAddress: p.WalletAddress,
	}
}

func ExistingUser() userdomain.User {
	return userdomain.User{ID: 1}
}
