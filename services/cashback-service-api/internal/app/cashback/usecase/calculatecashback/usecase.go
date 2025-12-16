package calculatecashback

import (
	"context"
	"errors"
	"log"

	"github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/domain"
	purchasedomain "github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/domain"
	userdomain "github.com/cashback-platform/services/cashback-service-api/internal/app/user/domain"
	"github.com/google/uuid"
)

const (
	DefaultCashbackPercent    = 5.0 // 5% default cashback
	EventTypeCashbackApproved = "cashback.approved"
)

type (
	// Repository interface for cashback operations
	Repository interface {
		Create(ctx context.Context, cashback domain.Cashback) (domain.Cashback, error)
		FindByPurchaseID(ctx context.Context, purchaseID uuid.UUID) (domain.Cashback, error)
	}

	// PurchaseRepository interface for purchase operations
	PurchaseRepository interface {
		FindByID(ctx context.Context, id uuid.UUID) (purchasedomain.Purchase, error)
	}

	// UserRepository interface for user operations
	UserRepository interface {
		FindByID(ctx context.Context, id uuid.UUID) (userdomain.User, error)
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
func (u UseCase) Execute(ctx context.Context, purchaseID uuid.UUID) (domain.Cashback, error) {
	existingCashback, err := u.repository.FindByPurchaseID(ctx, purchaseID)
	if err == nil {
		log.Printf("Cashback already exists for purchase %s", purchaseID)
		return existingCashback, ErrCashbackAlreadyExists
	}
	if !errors.Is(err, domain.ErrCashbackNotFound) {
		return domain.Cashback{}, err
	}

	purchase, err := u.purchaseRepository.FindByID(ctx, purchaseID)
	if err != nil {
		return domain.Cashback{}, ErrPurchaseNotFound
	}

	// Get user details (to validate and get wallet address)
	user, err := u.userRepository.FindByID(ctx, purchase.UserID)
	if err != nil {
		return domain.Cashback{}, ErrUserNotFound
	}

	// Calculate cashback
	cashback, err := domain.NewCashback(
		purchase.UserID,
		purchase.ID,
		purchase.Amount,
		DefaultCashbackPercent,
	)
	if err != nil {
		return domain.Cashback{}, err
	}

	// Approve cashback immediately (business rule: auto-approve)
	cashback.Approve()

	// Persist cashback
	cashback, err = u.repository.Create(ctx, cashback)
	if err != nil {
		return domain.Cashback{}, err
	}

	// Publish cashback.approved event for async minting
	event := CashbackApprovedEvent{
		CashbackID:      cashback.ID.String(),
		UserID:          cashback.UserID.String(),
		WalletAddress:   user.WalletAddress,
		PurchaseID:      cashback.PurchaseID.String(),
		Amount:          cashback.Amount,
		CashbackPercent: cashback.CashbackPercent,
	}

	if err := u.outboxPublisher.Publish(ctx, EventTypeCashbackApproved, event); err != nil {
		log.Printf("Failed to publish cashback.approved event: %v", err)
		return cashback, ErrFailedToPublishEvent
	}

	log.Printf("Cashback approved: %s for user %s, amount: %.2f",
		cashback.ID, cashback.UserID, cashback.Amount)

	return cashback, nil
}
