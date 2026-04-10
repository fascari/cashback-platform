package createuser

import (
	"strconv"

	"github.com/cashback-platform/kit/validator"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/user/domain"
)

type (
	InputPayload struct {
		ExternalID    string `json:"external_id"    validate:"required"`
		Email         string `json:"email"          validate:"required,email"`
		WalletAddress string `json:"wallet_address" validate:"required,min=20"`
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
	return validator.Validate(p)
}

func ToOutputPayload(user domain.User) OutputPayload {
	return OutputPayload{
		ID:            strconv.FormatInt(user.ID, 10),
		ExternalID:    user.ExternalID,
		Email:         user.Email,
		WalletAddress: user.WalletAddress,
		CreatedAt:     user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
