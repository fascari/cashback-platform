package testdata

import (
	"github.com/cashback-platform/services/cashback-service-api/internal/app/user/handler/createuser"
)

func ValidPayload() createuser.InputPayload {
	return createuser.InputPayload{
		ExternalID:    "ext-123",
		Email:         "user@example.com",
		WalletAddress: "0x1234567890abcdef1234",
	}
}

func InvalidEmailPayload() createuser.InputPayload {
	return createuser.InputPayload{ExternalID: "ext-123", Email: "not-an-email", WalletAddress: "0x1234567890abcdef1234"}
}

func ShortWalletPayload() createuser.InputPayload {
	return createuser.InputPayload{ExternalID: "ext-123", Email: "user@example.com", WalletAddress: "short"}
}
