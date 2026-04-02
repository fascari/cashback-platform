package calculatecashback

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/domain"
	purchasedomain "github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/domain"
	userdomain "github.com/cashback-platform/services/cashback-service-api/internal/app/user/domain"
)

const (
	DefaultCashbackPercent    = 5.0 // 5% default cashback
	EventTypeCashbackApproved = "cashback.approved"
)

type (
	// Repository interface for cashback operations
	Repository interface {
		Create(ctx context.Context, cashback domain.Cashback) (domain.Cashback, error)
		FindByPurchaseID(ctx context.Context, purchaseID int64) (domain.Cashback, error)
	}

	// PurchaseRepository interface for purchase operations
	PurchaseRepository interface {
		FindByID(ctx context.Context, id int64) (purchasedomain.Purchase, error)
	}

	// UserRepository interface for user operations
	UserRepository interface {
		FindByID(ctx context.Context, id int64) (userdomain.User, error)
	}

	// OutboxPublisher publishes events to the outbox
	OutboxPublisher interface {
		Publish(ctx context.Context, eventType string, payload any) error
	}

	// UseCase handles cashback calculation
	UseCase struct {
		repository         Repository
		purchaseRepository PurchaseRepository
		userRepository     UserRepository
		outboxPublisher    OutboxPublisher
	}

	// CashbackApprovedEvent represents the event published when cashback is approved
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
	outboxPublisher OutboxPublisher,
) UseCase {
	return UseCase{
		repository:         repository,
		purchaseRepository: purchaseRepository,
		userRepository:     userRepository,
		outboxPublisher:    outboxPublisher,
	}
}

// Execute calculates and creates cashback for a purchase
func (u UseCase) Execute(ctx context.Context, purchaseID int64) (domain.Cashback, error) {
	existingCashback, err := u.repository.FindByPurchaseID(ctx, purchaseID)
	if err == nil {
		log.Printf("Cashback already exists for purchase %d", purchaseID)
		return existingCashback, ErrCashbackAlreadyExists
	}
	if !errors.Is(err, domain.ErrCashbackNotFound) {
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

	cashback.Approve()

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

	if err := u.outboxPublisher.Publish(ctx, EventTypeCashbackApproved, event); err != nil {
		log.Printf("Failed to publish cashback.approved event: %v", err)
		return cashback, ErrFailedToPublishEvent
	}

	log.Printf("Cashback approved: %d for user %d, amount: %.2f",
		cashback.ID, cashback.UserID, cashback.Amount)

	return cashback, nil
}
