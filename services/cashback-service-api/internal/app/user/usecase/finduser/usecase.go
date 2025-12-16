package finduser

import (
	"context"

	"github.com/cashback-platform/services/cashback-service-api/internal/app/user/domain"
	"github.com/google/uuid"
)

type (
	Repository interface {
		FindByID(ctx context.Context, id uuid.UUID) (domain.User, error)
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

func (u UseCase) Execute(ctx context.Context, id uuid.UUID) (domain.User, error) {
	return u.repository.FindByID(ctx, id)
}
