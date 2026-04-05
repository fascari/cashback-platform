package processcashbackapproved

//go:generate mockery --all

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/cashback-platform/services/mint-consumer/internal/app/mint/domain"
	infragrpc "github.com/cashback-platform/services/mint-consumer/internal/infra/grpc"
)

type (
	// Repository defines the data access operations required by this use case.
	Repository interface {
		CreateMintRequest(ctx context.Context, req domain.MintRequest) (domain.MintRequest, error)
		MarkCompleted(ctx context.Context, id int64, txHash string, blockNumber int64) error
		MarkFailed(ctx context.Context, id int64, errorCode, errorMessage string, nextRetryAt *time.Time) error
		ExistsProcessedEvent(ctx context.Context, eventID uuid.UUID) (bool, error)
		CreateProcessedEvent(ctx context.Context, event domain.ProcessedEvent) error
	}

	// BlockchainClient is the gRPC gateway for token minting operations.
	BlockchainClient interface {
		MintToken(ctx context.Context, req infragrpc.MintTokenRequest) (infragrpc.MintResult, error)
	}

	// TransactionManager wraps multiple repository calls in a single DB transaction.
	TransactionManager interface {
		WithTransaction(ctx context.Context, fn func(context.Context) error) error
	}

	// UseCase handles an incoming cashback.approved NATS message.
	UseCase struct {
		repository         Repository
		blockchainClient   BlockchainClient
		transactionManager TransactionManager
	}

	cashbackApprovedEvent struct {
		EventID       string `json:"event_id"`
		CashbackID    string `json:"cashback_id"`
		UserID        string `json:"user_id"`
		WalletAddress string `json:"wallet_address"`
		PurchaseID    string `json:"purchase_id"`
		TokenAmount   string `json:"token_amount"`
		ChainID       string `json:"chain_id"`
	}

	parsedIDs struct {
		eventID        uuid.UUID
		cashbackID     int64
		userID         int64
		idempotencyKey uuid.UUID
	}
)

// NewUseCase returns a UseCase wired with its dependencies.
func NewUseCase(
	repository Repository,
	blockchainClient BlockchainClient,
	transactionManager TransactionManager,
) UseCase {
	return UseCase{
		repository:         repository,
		blockchainClient:   blockchainClient,
		transactionManager: transactionManager,
	}
}

// Execute processes a cashback.approved NATS message.
// Idempotency is enforced inside the transaction via processed_events.
// Returns nil on duplicate delivery so the caller can Ack without error.
func (u UseCase) Execute(ctx context.Context, payload []byte) error {
	event := new(cashbackApprovedEvent)
	if err := json.Unmarshal(payload, event); err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidPayload, err)
	}

	ids, err := parseIDs(event)
	if err != nil {
		return err
	}

	mintReq, err := u.createInTX(ctx, event, ids)
	if err != nil {
		return err
	}

	// mintReq.ID == 0: event was a duplicate — already processed.
	if mintReq.ID == 0 {
		return nil
	}

	result, err := u.blockchainClient.MintToken(ctx, infragrpc.MintTokenRequest{
		IdempotencyKey: ids.idempotencyKey.String(),
		WalletAddress:  event.WalletAddress,
		TokenAmount:    event.TokenAmount,
		ChainID:        event.ChainID,
	})
	if err != nil {
		nextRetry := retryAt(0)
		if markErr := u.repository.MarkFailed(ctx, mintReq.ID, "grpc_error", err.Error(), &nextRetry); markErr != nil {
			return fmt.Errorf("mark mint request failed after grpc error: %w", markErr)
		}
		return nil
	}

	return applyResult(ctx, u.repository, mintReq.ID, result, 0)
}

func (u UseCase) createInTX(ctx context.Context, event *cashbackApprovedEvent, ids parsedIDs) (domain.MintRequest, error) {
	var mintReq domain.MintRequest

	err := u.transactionManager.WithTransaction(ctx, func(txCtx context.Context) error {
		exists, err := u.repository.ExistsProcessedEvent(txCtx, ids.eventID)
		if err != nil {
			return fmt.Errorf("check processed event: %w", err)
		}
		if exists {
			return nil
		}

		mintReq, err = u.repository.CreateMintRequest(txCtx, domain.MintRequest{
			CashbackID:     ids.cashbackID,
			UserID:         ids.userID,
			WalletAddress:  event.WalletAddress,
			TokenAmount:    event.TokenAmount,
			IdempotencyKey: ids.idempotencyKey,
			Status:         domain.MintRequestStatusPending,
			MaxRetries:     5,
		})
		if err != nil {
			return fmt.Errorf("create mint request: %w", err)
		}

		if err := u.repository.CreateProcessedEvent(txCtx, domain.ProcessedEvent{
			EventID:   ids.eventID,
			EventType: "cashback.approved",
		}); err != nil {
			return fmt.Errorf("record processed event: %w", err)
		}

		return nil
	})

	return mintReq, err
}

func parseIDs(event *cashbackApprovedEvent) (parsedIDs, error) {
	eventID, err := uuid.Parse(event.EventID)
	if err != nil {
		return parsedIDs{}, fmt.Errorf("%w: event_id %q: %s", domain.ErrInvalidEventID, event.EventID, err)
	}

	cashbackID, err := strconv.ParseInt(event.CashbackID, 10, 64)
	if err != nil {
		return parsedIDs{}, fmt.Errorf("%w: cashback_id %q: %s", domain.ErrInvalidCashbackID, event.CashbackID, err)
	}

	userID, err := strconv.ParseInt(event.UserID, 10, 64)
	if err != nil {
		return parsedIDs{}, fmt.Errorf("%w: user_id %q: %s", domain.ErrInvalidUserID, event.UserID, err)
	}

	return parsedIDs{
		eventID:        eventID,
		cashbackID:     cashbackID,
		userID:         userID,
		idempotencyKey: uuid.NewSHA1(uuid.NameSpaceOID, []byte(event.EventID)),
	}, nil
}

func applyResult(ctx context.Context, repo Repository, mintReqID int64, result infragrpc.MintResult, retryCount int) error {
	if result.ErrorCode == "" {
		if err := repo.MarkCompleted(ctx, mintReqID, result.TransactionHash, result.BlockNumber); err != nil {
			return fmt.Errorf("mark mint request completed: %w", err)
		}
		return nil
	}

	var nextRetry *time.Time
	if result.Retryable {
		t := retryAt(retryCount)
		nextRetry = &t
	}

	if err := repo.MarkFailed(ctx, mintReqID, result.ErrorCode, result.ErrorMessage, nextRetry); err != nil {
		return fmt.Errorf("mark mint request failed: %w", err)
	}
	return nil
}

// retryAt calculates the next retry time using exponential backoff: 100ms * 2^retryCount.
func retryAt(retryCount int) time.Time {
	return time.Now().UTC().Add(100 * time.Millisecond * (1 << retryCount))
}
