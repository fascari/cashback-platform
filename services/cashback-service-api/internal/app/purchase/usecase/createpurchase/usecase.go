package createpurchase

import (
	"context"

	"github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/domain"
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

func (u UseCase) Execute(ctx context.Context, userID int64, amount float64, merchant string) (domain.Purchase, error) {
	if amount <= 0 {
		return domain.Purchase{}, ErrInvalidAmount
	}

	if userID == 0 {
		return domain.Purchase{}, ErrInvalidUserID
	}

	if merchant == "" {
		return domain.Purchase{}, ErrInvalidMerchant
	}

	purchase := domain.NewPurchase(userID, amount, merchant)
	return u.repository.Create(ctx, purchase)
}
