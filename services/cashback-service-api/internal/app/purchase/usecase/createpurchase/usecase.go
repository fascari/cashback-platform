package createpurchase

import (
	"context"

	"github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/domain"
	"github.com/google/uuid"
)

type (
	Repository interface {
		Create(ctx context.Context, purchase domain.Purchase) (domain.Purchase, error)
	}

	UseCase struct {
		repository Repository
	}
)

func New(repository Repository) UseCase {
	return UseCase{
		repository: repository,
	}
}

func (u UseCase) Execute(ctx context.Context, userID uuid.UUID, amount float64, merchant string) (domain.Purchase, error) {
	if amount <= 0 {
		return domain.Purchase{}, ErrInvalidAmount
	}

	if userID == uuid.Nil {
		return domain.Purchase{}, ErrInvalidUserID
	}

	if merchant == "" {
		return domain.Purchase{}, ErrInvalidMerchant
	}

	purchase := domain.NewPurchase(userID, amount, merchant)
	return u.repository.Create(ctx, purchase)
}
