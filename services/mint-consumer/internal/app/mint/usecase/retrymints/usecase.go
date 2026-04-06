package retrymints

//go:generate mockery --all

import (
	"context"
	"fmt"
	"time"

	"github.com/cashback-platform/services/mint-consumer/internal/app/mint/domain"
)

type (
	Repository interface {
		FindFailedRetryable(ctx context.Context, limit int) ([]domain.MintRequest, error)
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

func (u UseCase) Execute(ctx context.Context) error {
	requests, err := u.repository.FindFailedRetryable(ctx, 50)
	if err != nil {
		return fmt.Errorf("find retryable mint requests: %w", err)
	}

	for _, req := range requests {
		if err := u.retry(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

func (u UseCase) retry(ctx context.Context, req domain.MintRequest) error {
	result, err := u.blockchainClient.MintToken(ctx, domain.MintTokenRequest{
		IdempotencyKey: req.IdempotencyKey.String(),
		WalletAddress:  req.WalletAddress,
		TokenAmount:    req.TokenAmount,
	})
	if err != nil {
		if markErr := u.repository.MarkFailed(ctx, req.ID, "grpc_error", err.Error(), new(retryAt(req.RetryCount))); markErr != nil {
			return fmt.Errorf("mark mint request failed during retry: %w", markErr)
		}
		return nil
	}

	return applyResult(ctx, u.repository, req.ID, result, req.RetryCount)
}

func applyResult(ctx context.Context, repo Repository, mintReqID int64, result domain.MintResult, retryCount int) error {
	if result.ErrorCode == "" {
		if err := repo.MarkCompleted(ctx, mintReqID, result.TransactionHash, result.BlockNumber); err != nil {
			return fmt.Errorf("mark mint request completed: %w", err)
		}
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
