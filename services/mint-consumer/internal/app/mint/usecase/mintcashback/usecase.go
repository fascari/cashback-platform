package mintcashback

//go:generate mockery --all

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/cashback-platform/kit/events"
	"github.com/cashback-platform/kit/logger"
	"github.com/cashback-platform/services/mint-consumer/internal/app/mint/domain"
)

type (
	Repository interface {
		CreateMintRequestIdempotent(ctx context.Context, req domain.MintRequest, eventID uuid.UUID, eventType string) (domain.MintRequest, bool, error)
		MarkCompleted(ctx context.Context, id int64, txHash string, blockNumber int64) error
		MarkFailed(ctx context.Context, id int64, errorCode, errorMessage string, nextRetryAt *time.Time) error
	}

	BlockchainClient interface {
		MintToken(ctx context.Context, req domain.MintTokenRequest) (domain.MintResult, error)
	}

	UseCase struct {
		repository       Repository
		blockchainClient BlockchainClient
	}
)

func NewUseCase(repository Repository, blockchainClient BlockchainClient) UseCase {
	return UseCase{repository: repository, blockchainClient: blockchainClient}
}

func (u UseCase) Execute(ctx context.Context, input Input) error {
	idempotencyKey := uuid.NewSHA1(uuid.NameSpaceOID, []byte(input.EventID.String()))

	mintReq, isNew, err := u.repository.CreateMintRequestIdempotent(ctx, domain.MintRequest{
		CashbackID:     input.CashbackID,
		UserID:         input.UserID,
		WalletAddress:  input.WalletAddress,
		TokenAmount:    input.TokenAmount,
		IdempotencyKey: idempotencyKey,
		Status:         domain.MintRequestStatusPending,
		MaxRetries:     5,
	}, input.EventID, events.CashbackApproved)
	if err != nil {
		return fmt.Errorf("create mint request idempotent: %w", err)
	}

	// mintReq.ID == 0: event was a duplicate — already processed.
	if !isNew {
		logger.Info("mint skipped: duplicate event", "event_id", input.EventID)
		return nil
	}

	logger.Info("minting cashback", "cashback_id", input.CashbackID, "wallet", input.WalletAddress, "amount", input.TokenAmount)

	result, err := u.blockchainClient.MintToken(ctx, domain.MintTokenRequest{
		IdempotencyKey: idempotencyKey.String(),
		WalletAddress:  input.WalletAddress,
		TokenAmount:    input.TokenAmount,
		ChainID:        input.ChainID,
	})
	if err != nil {
		if markErr := u.repository.MarkFailed(ctx, mintReq.ID, "grpc_error", err.Error(), new(retryAt(0))); markErr != nil {
			return fmt.Errorf("mark mint request failed after grpc error: %w", markErr)
		}
		return nil
	}

	return applyResult(ctx, u.repository, mintReq.ID, result, 0)
}

func applyResult(ctx context.Context, repo Repository, mintReqID int64, result domain.MintResult, retryCount int) error {
	if result.ErrorCode == "" {
		if err := repo.MarkCompleted(ctx, mintReqID, result.TransactionHash, result.BlockNumber); err != nil {
			return fmt.Errorf("mark mint request completed: %w", err)
		}
		logger.Info("mint completed", "mint_id", mintReqID, "tx_hash", result.TransactionHash, "block", result.BlockNumber)
		return nil
	}

	var nextRetry *time.Time
	if result.Retryable {
		nextRetry = new(retryAt(retryCount))
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
