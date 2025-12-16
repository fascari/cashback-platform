package findusercashback

import (
	"context"

	"github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/domain"
	"github.com/google/uuid"
)

type (
	Repository interface {
		FindByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Cashback, error)
		TotalByUserID(ctx context.Context, userID uuid.UUID) (float64, error)
	}

	UseCase struct {
		repository Repository
	}

	UserCashbackSummary struct {
		UserID         uuid.UUID
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

func (u UseCase) Execute(ctx context.Context, userID uuid.UUID) (UserCashbackSummary, error) {
	if userID == uuid.Nil {
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
