package createuser

import (
	"github.com/cashback-platform/services/cashback-service-api/internal/app/user/domain"
	"github.com/cashback-platform/services/cashback-service-api/pkg/validator"
)

type (
	InputPayload struct {
		ExternalID    string `json:"external_id"`
		Email         string `json:"email"`
		WalletAddress string `json:"wallet_address"`
	}

	OutputPayload struct {
		ID            string `json:"id"`
		ExternalID    string `json:"external_id"`
		Email         string `json:"email"`
		WalletAddress string `json:"wallet_address"`
		CreatedAt     string `json:"created_at"`
	}
)

func (p InputPayload) Validate() error {
	if p.ExternalID == "" {
		return validator.ErrRequired
	}
	if err := validator.ValidateEmail(p.Email); err != nil {
		return err
	}
	return validator.ValidateWalletAddress(p.WalletAddress)
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
