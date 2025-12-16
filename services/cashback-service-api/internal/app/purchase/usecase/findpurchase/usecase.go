package findpurchase

import (
	"context"

	"github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/domain"
	"github.com/google/uuid"
)

type (
	Repository interface {
		FindByID(ctx context.Context, id uuid.UUID) (domain.Purchase, error)
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

func (u UseCase) Execute(ctx context.Context, id uuid.UUID) (domain.Purchase, error) {
	return u.repository.FindByID(ctx, id)
}
