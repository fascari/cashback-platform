package createuser

import (
	"context"
	"time"

	"github.com/cashback-platform/services/cashback-service-api/internal/app/user/domain"
	"github.com/google/uuid"
)

type (
	Repository interface {
		Create(ctx context.Context, user domain.User) (domain.User, error)
		FindByEmail(ctx context.Context, email string) (domain.User, error)
		FindByExternalID(ctx context.Context, externalID string) (domain.User, error)
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

func (u UseCase) Execute(ctx context.Context, externalID, email, walletAddress string) (domain.User, error) {
	if existingUser, _ := u.repository.FindByEmail(ctx, email); existingUser.ID != uuid.Nil {
		return domain.User{}, ErrUserAlreadyExists
	}

	if existingUser, _ := u.repository.FindByExternalID(ctx, externalID); existingUser.ID != uuid.Nil {
		return domain.User{}, ErrUserAlreadyExists
	}

	user := domain.User{
		ID:            uuid.New(),
		ExternalID:    externalID,
		Email:         email,
		WalletAddress: walletAddress,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}

	return u.repository.Create(ctx, user)
}
