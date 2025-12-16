package finduser

import (
	"github.com/cashback-platform/services/cashback-service-api/internal/app/user/domain"
)

type OutputPayload struct {
	ID            string `json:"id"`
	ExternalID    string `json:"external_id"`
	Email         string `json:"email"`
	WalletAddress string `json:"wallet_address"`
	CreatedAt     string `json:"created_at"`
}

func ToOutputPayload(user domain.User) OutputPayload {
	return OutputPayload{
		ID:            user.ID.String(),
		ExternalID:    user.ExternalID,
		Email:         user.Email,
		WalletAddress: user.WalletAddress,
		CreatedAt:     user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
