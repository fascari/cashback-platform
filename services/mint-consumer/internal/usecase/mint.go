package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/cashback-platform/services/mint-consumer/internal/domain"
)

type (
	MintRequestRepository interface {
		Create(ctx context.Context, req domain.MintRequest) (domain.MintRequest, error)
		FindFailedRetryable(ctx context.Context, limit int) ([]domain.MintRequest, error)
		MarkCompleted(ctx context.Context, id int64, txHash string, blockNumber int64) error
		MarkFailed(ctx context.Context, id int64, errorCode, errorMessage string, nextRetryAt *time.Time) error
	}

	ProcessedEventRepository interface {
		Exists(ctx context.Context, eventID uuid.UUID) (bool, error)
		Create(ctx context.Context, event domain.ProcessedEvent) error
	}

	MintTokenRequest struct {
		IdempotencyKey string
		WalletAddress  string
		TokenAmount    string
		ChainID        string
	}

	MintResult struct {
		TransactionHash string
		BlockNumber     int64
		ErrorCode       string
		ErrorMessage    string
		Retryable       bool
	}

	BlockchainClient interface {
		MintToken(ctx context.Context, req MintTokenRequest) (MintResult, error)
	}

	TransactionManager interface {
		WithTransaction(ctx context.Context, fn func(context.Context) error) error
	}

	MintUsecase struct {
		mintRequestRepo    MintRequestRepository
		processedEventRepo ProcessedEventRepository
		blockchainClient   BlockchainClient
		txManager          TransactionManager
	}

	cashbackEventIDs struct {
		eventID        uuid.UUID
		cashbackID     int64
		userID         int64
		idempotencyKey uuid.UUID
	}
)

func NewMint(
	mintRequestRepo MintRequestRepository,
	processedEventRepo ProcessedEventRepository,
	blockchainClient BlockchainClient,
	txManager TransactionManager,
) MintUsecase {
	return MintUsecase{
		mintRequestRepo:    mintRequestRepo,
		processedEventRepo: processedEventRepo,
		blockchainClient:   blockchainClient,
		txManager:          txManager,
	}
}

// ProcessCashbackApproved handles a cashback.approved NATS message.
// Idempotency: duplicate deliveries are detected inside the transaction via processedEvent;
// gRPC is never called for duplicates. Returns nil on success so the caller can Ack.
func (u MintUsecase) ProcessCashbackApproved(ctx context.Context, payload []byte) error {
	event := new(CashbackApprovedEvent)
	if err := json.Unmarshal(payload, event); err != nil {
		return fmt.Errorf("unmarshal cashback approved event: %w", err)
	}

	ids, err := parseCashbackIDs(event)
	if err != nil {
		return err
	}

	mintReq, err := u.createMintInTX(ctx, event, ids)
	if err != nil {
		return err
	}

	// mintReq.ID == 0: event was a duplicate — already processed.
	if mintReq.ID == 0 {
		return nil
	}

	result, err := u.blockchainClient.MintToken(ctx, MintTokenRequest{
		IdempotencyKey: ids.idempotencyKey.String(),
		WalletAddress:  event.WalletAddress,
		TokenAmount:    event.TokenAmount,
		ChainID:        event.ChainID,
	})
	if err != nil {
		nextRetry := retryAt(0)
		if markErr := u.mintRequestRepo.MarkFailed(ctx, mintReq.ID, "grpc_error", err.Error(), &nextRetry); markErr != nil {
			return fmt.Errorf("mark mint request failed after grpc error: %w", markErr)
		}
		return nil
	}

	return u.applyMintResult(ctx, mintReq.ID, result, 0)
}

// RetryFailedMints retries mint requests that failed with a retryable error and are due for retry.
func (u MintUsecase) RetryFailedMints(ctx context.Context) error {
	requests, err := u.mintRequestRepo.FindFailedRetryable(ctx, 50)
	if err != nil {
		return fmt.Errorf("find failed retryable mint requests: %w", err)
	}

	for _, req := range requests {
		if err := u.retryMintRequest(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

func (u MintUsecase) createMintInTX(ctx context.Context, event *CashbackApprovedEvent, ids cashbackEventIDs) (domain.MintRequest, error) {
	var mintReq domain.MintRequest

	txErr := u.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		exists, err := u.processedEventRepo.Exists(txCtx, ids.eventID)
		if err != nil {
			return fmt.Errorf("check processed event: %w", err)
		}
		if exists {
			return nil
		}

		mintReq, err = u.mintRequestRepo.Create(txCtx, domain.MintRequest{
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

		if err := u.processedEventRepo.Create(txCtx, domain.ProcessedEvent{
			EventID:   ids.eventID,
			EventType: "cashback.approved",
		}); err != nil {
			return fmt.Errorf("record processed event: %w", err)
		}

		return nil
	})

	return mintReq, txErr
}

func (u MintUsecase) retryMintRequest(ctx context.Context, req domain.MintRequest) error {
	result, err := u.blockchainClient.MintToken(ctx, MintTokenRequest{
		IdempotencyKey: req.IdempotencyKey.String(),
		WalletAddress:  req.WalletAddress,
		TokenAmount:    req.TokenAmount,
	})
	if err != nil {
		nextRetry := retryAt(req.RetryCount)
		if markErr := u.mintRequestRepo.MarkFailed(ctx, req.ID, "grpc_error", err.Error(), &nextRetry); markErr != nil {
			return fmt.Errorf("mark mint request failed during retry: %w", markErr)
		}
		return nil
	}

	return u.applyMintResult(ctx, req.ID, result, req.RetryCount)
}

func (u MintUsecase) applyMintResult(ctx context.Context, mintReqID int64, result MintResult, retryCount int) error {
	if result.ErrorCode == "" {
		if err := u.mintRequestRepo.MarkCompleted(ctx, mintReqID, result.TransactionHash, result.BlockNumber); err != nil {
			return fmt.Errorf("mark mint request completed: %w", err)
		}
		return nil
	}

	var nextRetry *time.Time
	if result.Retryable {
		t := retryAt(retryCount)
		nextRetry = &t
	}

	if err := u.mintRequestRepo.MarkFailed(ctx, mintReqID, result.ErrorCode, result.ErrorMessage, nextRetry); err != nil {
		return fmt.Errorf("mark mint request failed: %w", err)
	}
	return nil
}

func parseCashbackIDs(event *CashbackApprovedEvent) (cashbackEventIDs, error) {
	eventID, err := uuid.Parse(event.EventID)
	if err != nil {
		return cashbackEventIDs{}, fmt.Errorf("parse event id %q: %w", event.EventID, err)
	}

	cashbackID, err := strconv.ParseInt(event.CashbackID, 10, 64)
	if err != nil {
		return cashbackEventIDs{}, fmt.Errorf("parse cashback id %q: %w", event.CashbackID, err)
	}

	userID, err := strconv.ParseInt(event.UserID, 10, 64)
	if err != nil {
		return cashbackEventIDs{}, fmt.Errorf("parse user id %q: %w", event.UserID, err)
	}

	return cashbackEventIDs{
		eventID:        eventID,
		cashbackID:     cashbackID,
		userID:         userID,
		idempotencyKey: uuid.NewSHA1(uuid.NameSpaceOID, []byte(event.EventID)),
	}, nil
}

// retryAt calculates the next retry time using exponential backoff: 100ms * 2^retryCount.
func retryAt(retryCount int) time.Time {
	backoff := 100 * time.Millisecond * (1 << retryCount)
	return time.Now().UTC().Add(backoff)
}
