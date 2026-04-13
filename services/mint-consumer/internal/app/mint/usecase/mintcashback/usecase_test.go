package mintcashback_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/cashback-platform/kit/events"
	"github.com/cashback-platform/services/mint-consumer/internal/app/mint/domain"
	"github.com/cashback-platform/services/mint-consumer/internal/app/mint/usecase/mintcashback"
	"github.com/cashback-platform/services/mint-consumer/internal/app/mint/usecase/mintcashback/mocks"
	"github.com/cashback-platform/services/mint-consumer/internal/app/mint/usecase/mintcashback/testdata"
)

func TestExecute_ShouldCreateMintRequestAndMarkCompletedWhenEventIsNew(t *testing.T) {
	repo := mocks.NewRepository(t)
	repo.EXPECT().CreateMintRequestIdempotent(mock.Anything, mock.Anything, testdata.EventID, events.CashbackApproved).
		Return(testdata.PendingMintRequest(), true, nil)
	repo.EXPECT().MarkCompleted(mock.Anything, testdata.PendingMintRequest().ID, "0xhash", int64(12345)).Return(nil)

	bc := mocks.NewBlockchainClient(t)
	bc.EXPECT().MintToken(mock.Anything, mock.Anything).Return(testdata.SuccessfulMintResult(), nil)

	uc := mintcashback.NewUseCase(repo, bc)
	err := uc.Execute(context.Background(), testdata.NewInput())

	require.NoError(t, err)
}

func TestExecute_ShouldSkipGRPCCallWhenEventIsDuplicate(t *testing.T) {
	repo := mocks.NewRepository(t)
	repo.EXPECT().CreateMintRequestIdempotent(mock.Anything, mock.Anything, testdata.EventID, events.CashbackApproved).
		Return(domain.MintRequest{}, false, nil)

	bc := mocks.NewBlockchainClient(t)

	uc := mintcashback.NewUseCase(repo, bc)
	err := uc.Execute(context.Background(), testdata.NewInput())

	require.NoError(t, err)
}

func TestExecute_ShouldMarkFailedWithRetryWhenGRPCReturnsRetryableResult(t *testing.T) {
	repo := mocks.NewRepository(t)
	repo.EXPECT().CreateMintRequestIdempotent(mock.Anything, mock.Anything, testdata.EventID, events.CashbackApproved).
		Return(testdata.PendingMintRequest(), true, nil)
	repo.EXPECT().MarkFailed(mock.Anything, testdata.PendingMintRequest().ID, "timeout", "rpc timeout", mock.MatchedBy(func(t *time.Time) bool {
		return t != nil
	})).Return(nil)

	bc := mocks.NewBlockchainClient(t)
	bc.EXPECT().MintToken(mock.Anything, mock.Anything).Return(testdata.RetryableMintResult(), nil)

	uc := mintcashback.NewUseCase(repo, bc)
	err := uc.Execute(context.Background(), testdata.NewInput())

	require.NoError(t, err)
}

func TestExecute_ShouldMarkFailedWithoutRetryWhenGRPCReturnsPermanentResult(t *testing.T) {
	repo := mocks.NewRepository(t)
	repo.EXPECT().CreateMintRequestIdempotent(mock.Anything, mock.Anything, testdata.EventID, events.CashbackApproved).
		Return(testdata.PendingMintRequest(), true, nil)
	repo.EXPECT().MarkFailed(mock.Anything, testdata.PendingMintRequest().ID, "invalid_address", "bad address", (*time.Time)(nil)).Return(nil)

	bc := mocks.NewBlockchainClient(t)
	bc.EXPECT().MintToken(mock.Anything, mock.Anything).Return(testdata.PermanentFailMintResult(), nil)

	uc := mintcashback.NewUseCase(repo, bc)
	err := uc.Execute(context.Background(), testdata.NewInput())

	require.NoError(t, err)
}

func TestExecute_ShouldReturnErrorWhenCreateMintRequestIdempotentFails(t *testing.T) {
	repo := mocks.NewRepository(t)
	repo.EXPECT().CreateMintRequestIdempotent(mock.Anything, mock.Anything, testdata.EventID, events.CashbackApproved).
		Return(domain.MintRequest{}, false, errors.New("db error"))

	bc := mocks.NewBlockchainClient(t)

	uc := mintcashback.NewUseCase(repo, bc)
	err := uc.Execute(context.Background(), testdata.NewInput())

	require.Error(t, err)
}
