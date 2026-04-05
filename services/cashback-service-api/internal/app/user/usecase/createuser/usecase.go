package createuser

//go:generate mockery --all

import (
	"context"
	"time"

	"github.com/cashback-platform/services/cashback-service-api/internal/app/user/domain"
	"github.com/cashback-platform/kit/apperror"
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
	return UseCase{repository: repository}
}

func (u UseCase) Execute(ctx context.Context, externalID, email, walletAddress string) (domain.User, error) {
	_, err := u.repository.FindByEmail(ctx, email)
	if err == nil {
		return domain.User{}, domain.ErrUserAlreadyExists
	}
	if !apperror.As(err, domain.ErrCodeUserNotFound) {
		return domain.User{}, err
	}

	_, err = u.repository.FindByExternalID(ctx, externalID)
	if err == nil {
		return domain.User{}, domain.ErrUserAlreadyExists
	}
	if !apperror.As(err, domain.ErrCodeUserNotFound) {
		return domain.User{}, err
	}

	user := domain.User{
		ExternalID:    externalID,
		Email:         email,
		WalletAddress: walletAddress,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}

	return u.repository.Create(ctx, user)
}
