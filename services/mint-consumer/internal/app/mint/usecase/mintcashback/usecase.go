package mintcashback

//go:generate mockery --all

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/cashback-platform/services/mint-consumer/internal/app/mint/domain"
)

type (
	Repository interface {
		CreateMintRequest(ctx context.Context, req domain.MintRequest) (domain.MintRequest, error)
		MarkCompleted(ctx context.Context, id int64, txHash string, blockNumber int64) error
		MarkFailed(ctx context.Context, id int64, errorCode, errorMessage string, nextRetryAt *time.Time) error
		ExistsProcessedEvent(ctx context.Context, eventID uuid.UUID) (bool, error)
		CreateProcessedEvent(ctx context.Context, eventID uuid.UUID, eventType string) error
	}

	BlockchainClient interface {
		MintToken(ctx context.Context, req domain.MintTokenRequest) (domain.MintResult, error)
	}

	TransactionManager interface {
		WithTransaction(ctx context.Context, fn func(context.Context) error) error
	}

	UseCase struct {
		repository         Repository
		blockchainClient   BlockchainClient
		transactionManager TransactionManager
	}

	Input struct {
		EventID       uuid.UUID
		CashbackID    int64
		UserID        int64
		WalletAddress string
		PurchaseID    string
		TokenAmount   string
		ChainID       string
	}
)

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

func (u UseCase) Execute(ctx context.Context, input Input) error {
	idempotencyKey := uuid.NewSHA1(uuid.NameSpaceOID, []byte(input.EventID.String()))

	mintReq, err := u.createInTX(ctx, input, idempotencyKey)
	if err != nil {
		return err
	}

	// mintReq.ID == 0: event was a duplicate — already processed.
	if mintReq.ID == 0 {
		return nil
	}

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

func (u UseCase) createInTX(ctx context.Context, input Input, idempotencyKey uuid.UUID) (domain.MintRequest, error) {
	var mintReq domain.MintRequest

	err := u.transactionManager.WithTransaction(ctx, func(txCtx context.Context) error {
		exists, err := u.repository.ExistsProcessedEvent(txCtx, input.EventID)
		if err != nil {
			return fmt.Errorf("check processed event: %w", err)
		}
		if exists {
			return nil
		}

		mintReq, err = u.repository.CreateMintRequest(txCtx, domain.MintRequest{
			CashbackID:     input.CashbackID,
			UserID:         input.UserID,
			WalletAddress:  input.WalletAddress,
			TokenAmount:    input.TokenAmount,
			IdempotencyKey: idempotencyKey,
			Status:         domain.MintRequestStatusPending,
			MaxRetries:     5,
		})
		if err != nil {
			return fmt.Errorf("create mint request: %w", err)
		}

		if err := u.repository.CreateProcessedEvent(txCtx, input.EventID, "cashback.approved"); err != nil {
			return fmt.Errorf("record processed event: %w", err)
		}

		return nil
	})

	return mintReq, err
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
