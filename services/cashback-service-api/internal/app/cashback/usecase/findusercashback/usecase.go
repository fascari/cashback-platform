package findusercashback

import (
	"context"

	"github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/domain"
)

type (
	Repository interface {
		FindByUserID(ctx context.Context, userID int64) ([]domain.Cashback, error)
		TotalByUserID(ctx context.Context, userID int64) (float64, error)
	}

	UseCase struct {
		repository Repository
	}

	UserCashbackSummary struct {
		UserID         int64
		Cashbacks      []domain.Cashback
		TotalMinted    float64
		TotalCashbacks int
	}
)

func New(repository Repository) UseCase {
	return UseCase{
		repository: repository,
	}
}

func (u UseCase) Execute(ctx context.Context, userID int64) (UserCashbackSummary, error) {
	if userID == 0 {
		return UserCashbackSummary{}, domain.ErrInvalidUserID
	}

	cashbacks, err := u.repository.FindByUserID(ctx, userID)
	if err != nil {
		return UserCashbackSummary{}, err
	}

	totalMinted, err := u.repository.TotalByUserID(ctx, userID)
	if err != nil {
		return UserCashbackSummary{}, err
	}

	return UserCashbackSummary{
		UserID:         userID,
		Cashbacks:      cashbacks,
		TotalMinted:    totalMinted,
		TotalCashbacks: len(cashbacks),
	}, nil
}
