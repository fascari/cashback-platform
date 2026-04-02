package calculatecashback

//go:generate mockery --all

import (
	"context"
	"strconv"

	"github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/domain"
	"github.com/cashback-platform/services/cashback-service-api/pkg/apperror"
)

const (
	DefaultCashbackPercent    = 5.0
	EventTypeCashbackApproved = "cashback.approved"
)

type (
	Repository interface {
		Create(ctx context.Context, cashback domain.Cashback) (domain.Cashback, error)
		FindByPurchaseID(ctx context.Context, purchaseID int64) (domain.Cashback, error)
	}

	Purchase struct {
		ID     int64
		UserID int64
		Amount float64
	}

	User struct {
		WalletAddress string
	}

	PurchaseRepository interface {
		FindByID(ctx context.Context, id int64) (Purchase, error)
	}

	UserRepository interface {
		FindByID(ctx context.Context, id int64) (User, error)
	}

	EventPublisher interface {
		Publish(ctx context.Context, eventType string, aggregateID int64, payload any) error
	}

	UseCase struct {
		repository         Repository
		purchaseRepository PurchaseRepository
		userRepository     UserRepository
		eventPublisher     EventPublisher
	}

	CashbackApprovedEvent struct {
		CashbackID      string  `json:"cashback_id"`
		UserID          string  `json:"user_id"`
		WalletAddress   string  `json:"wallet_address"`
		PurchaseID      string  `json:"purchase_id"`
		Amount          float64 `json:"amount"`
		CashbackPercent float64 `json:"cashback_percent"`
	}
)

func New(
	repository Repository,
	purchaseRepository PurchaseRepository,
	userRepository UserRepository,
	eventPublisher EventPublisher,
) UseCase {
	return UseCase{
		repository:         repository,
		purchaseRepository: purchaseRepository,
		userRepository:     userRepository,
		eventPublisher:     eventPublisher,
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

	cashback, err = u.repository.Create(ctx, cashback)
	if err != nil {
		return domain.Cashback{}, err
	}

	event := CashbackApprovedEvent{
		CashbackID:      strconv.FormatInt(cashback.ID, 10),
		UserID:          strconv.FormatInt(cashback.UserID, 10),
		WalletAddress:   user.WalletAddress,
		PurchaseID:      strconv.FormatInt(cashback.PurchaseID, 10),
		Amount:          cashback.Amount,
		CashbackPercent: cashback.CashbackPercent,
	}

	if err := u.eventPublisher.Publish(ctx, EventTypeCashbackApproved, cashback.ID, event); err != nil {
		return cashback, ErrFailedToPublishEvent
	}

	return cashback, nil
}
