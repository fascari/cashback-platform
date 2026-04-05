package calculatecashback

//go:generate mockery --all

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/domain"
	"github.com/cashback-platform/services/cashback-service-api/pkg/apperror"
)

const (
	DefaultCashbackPercent    = 5.0
	EventTypeCashbackApproved = "cashback.approved"

	defaultChainID = "ethereum"
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

	TransactionManager interface {
		WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
	}

	UseCase struct {
		repository         Repository
		purchaseRepository PurchaseRepository
		userRepository     UserRepository
		eventPublisher     EventPublisher
		transactionManager TransactionManager
	}

	CashbackApprovedEvent struct {
		EventID       string `json:"event_id"`
		CashbackID    string `json:"cashback_id"`
		UserID        string `json:"user_id"`
		WalletAddress string `json:"wallet_address"`
		PurchaseID    string `json:"purchase_id"`
		TokenAmount   string `json:"token_amount"`
		ChainID       string `json:"chain_id"`
	}
)

func New(
	repository Repository,
	purchaseRepository PurchaseRepository,
	userRepository UserRepository,
	eventPublisher EventPublisher,
	transactionManager TransactionManager,
) UseCase {
	return UseCase{
		repository:         repository,
		purchaseRepository: purchaseRepository,
		userRepository:     userRepository,
		eventPublisher:     eventPublisher,
		transactionManager: transactionManager,
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

	var created domain.Cashback
	if err := u.transactionManager.WithTransaction(ctx, func(txCtx context.Context) error {
		var txErr error
		created, txErr = u.repository.Create(txCtx, cashback)
		if txErr != nil {
			return txErr
		}

		event := CashbackApprovedEvent{
			EventID:       strconv.FormatInt(created.ID, 10),
			CashbackID:    strconv.FormatInt(created.ID, 10),
			UserID:        strconv.FormatInt(created.UserID, 10),
			WalletAddress: user.WalletAddress,
			PurchaseID:    strconv.FormatInt(created.PurchaseID, 10),
			TokenAmount:   fmt.Sprintf("%.0f", created.Amount*1e18),
			ChainID:       defaultChainID,
		}

		if err := u.eventPublisher.Publish(txCtx, EventTypeCashbackApproved, created.ID, event); err != nil {
			return ErrFailedToPublishEvent
		}
		return nil
	}); err != nil {
		return created, err
	}

	return created, nil
}
