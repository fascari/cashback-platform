package calculatecashback

//go:generate mockery --all

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cashback-platform/kit/apperror"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/domain"
	purchasedomain "github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/domain"
	userdomain "github.com/cashback-platform/services/cashback-service-api/internal/app/user/domain"
	"github.com/google/uuid"
)

const (
	DefaultCashbackPercent = 5.0

	defaultChainID = "ethereum"
)

type (
	Repository interface {
		FindByPurchaseID(ctx context.Context, purchaseID int64) (domain.Cashback, error)
		CreateWithEvent(ctx context.Context, cashback domain.Cashback, buildPayload func(domain.Cashback) any) (domain.Cashback, error)
	}

	PurchaseRepository interface {
		FindByID(ctx context.Context, id int64) (purchasedomain.Purchase, error)
	}

	UserRepository interface {
		FindByID(ctx context.Context, id int64) (userdomain.User, error)
	}

	UseCase struct {
		repository         Repository
		purchaseRepository PurchaseRepository
		userRepository     UserRepository
	}
)

func New(
	repository Repository,
	purchaseRepository PurchaseRepository,
	userRepository UserRepository,
) UseCase {
	return UseCase{
		repository:         repository,
		purchaseRepository: purchaseRepository,
		userRepository:     userRepository,
	}
}

func (u UseCase) Execute(ctx context.Context, purchaseID int64) (domain.Cashback, error) {
	existingCashback, err := u.repository.FindByPurchaseID(ctx, purchaseID)
	if err == nil {
		return existingCashback, domain.ErrCashbackAlreadyExists
	}
	if !apperror.As(err, domain.ErrCodeCashbackNotFound) {
		return domain.Cashback{}, err
	}

	purchase, err := u.purchaseRepository.FindByID(ctx, purchaseID)
	if err != nil {
		return domain.Cashback{}, ErrPurchaseNotFound
	}

	user, err := u.userRepository.FindByID(ctx, purchase.UserID)
	if err != nil {
		return domain.Cashback{}, ErrUserNotFound
	}

	cashback, err := domain.NewCashback(
		purchase.UserID,
		purchase.ID,
		purchase.Amount,
		DefaultCashbackPercent,
	)
	if err != nil {
		return domain.Cashback{}, err
	}

	cashback = cashback.Approve()

	return u.repository.CreateWithEvent(ctx, cashback, func(created domain.Cashback) any {
		purchaseID := ""
		if created.PurchaseID != nil {
			purchaseID = strconv.FormatInt(*created.PurchaseID, 10)
		}
		return CashbackApprovedEvent{
			EventID:       uuid.New().String(),
			CashbackID:    strconv.FormatInt(created.ID, 10),
			UserID:        strconv.FormatInt(created.UserID, 10),
			WalletAddress: user.WalletAddress,
			PurchaseID:    purchaseID,
			TokenAmount:   fmt.Sprintf("%.0f", created.Amount*1e18),
			ChainID:       defaultChainID,
		}
	})
}
