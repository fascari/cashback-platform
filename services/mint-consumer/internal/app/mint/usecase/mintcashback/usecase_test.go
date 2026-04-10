package mintcashback_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/cashback-platform/services/mint-consumer/internal/app/mint/usecase/mintcashback"
	"github.com/cashback-platform/services/mint-consumer/internal/app/mint/usecase/mintcashback/mocks"
	"github.com/cashback-platform/services/mint-consumer/internal/app/mint/usecase/mintcashback/testdata"
)

func txManagerPassthrough(t *testing.T) *mocks.TransactionManager {
	t.Helper()
	tm := mocks.NewTransactionManager(t)
	tm.EXPECT().WithTransaction(mock.Anything, mock.Anything).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})
	return tm
}

func TestExecute_ShouldCreateMintRequestAndMarkCompletedWhenEventIsNew(t *testing.T) {
	repo := mocks.NewRepository(t)
	repo.EXPECT().ExistsProcessedEvent(mock.Anything, testdata.EventID).Return(false, nil)
	repo.EXPECT().CreateMintRequest(mock.Anything, mock.Anything).Return(testdata.PendingMintRequest(), nil)
	repo.EXPECT().CreateProcessedEvent(mock.Anything, testdata.EventID, "cashback.approved").Return(nil)
	repo.EXPECT().MarkCompleted(mock.Anything, testdata.PendingMintRequest().ID, "0xhash", int64(12345)).Return(nil)

	bc := mocks.NewBlockchainClient(t)
	bc.EXPECT().MintToken(mock.Anything, mock.Anything).Return(testdata.SuccessfulMintResult(), nil)

	uc := mintcashback.NewUseCase(repo, bc, txManagerPassthrough(t))
	err := uc.Execute(context.Background(), testdata.NewInput())

	require.NoError(t, err)
}

func TestExecute_ShouldSkipGRPCCallWhenEventIsDuplicate(t *testing.T) {
	repo := mocks.NewRepository(t)
	repo.EXPECT().ExistsProcessedEvent(mock.Anything, testdata.EventID).Return(true, nil)

	bc := mocks.NewBlockchainClient(t)

	uc := mintcashback.NewUseCase(repo, bc, txManagerPassthrough(t))
	err := uc.Execute(context.Background(), testdata.NewInput())

	require.NoError(t, err)
}

func TestExecute_ShouldMarkFailedWithRetryWhenGRPCReturnsRetryableResult(t *testing.T) {
	repo := mocks.NewRepository(t)
	repo.EXPECT().ExistsProcessedEvent(mock.Anything, testdata.EventID).Return(false, nil)
	repo.EXPECT().CreateMintRequest(mock.Anything, mock.Anything).Return(testdata.PendingMintRequest(), nil)
	repo.EXPECT().CreateProcessedEvent(mock.Anything, testdata.EventID, "cashback.approved").Return(nil)
	repo.EXPECT().MarkFailed(mock.Anything, testdata.PendingMintRequest().ID, "timeout", "rpc timeout", mock.MatchedBy(func(t *time.Time) bool {
		return t != nil
	})).Return(nil)

	result := testdata.RetryableMintResult()
	bc := mocks.NewBlockchainClient(t)
	bc.EXPECT().MintToken(mock.Anything, mock.Anything).Return(result, nil)

	uc := mintcashback.NewUseCase(repo, bc, txManagerPassthrough(t))
	err := uc.Execute(context.Background(), testdata.NewInput())

	require.NoError(t, err)
}

func TestExecute_ShouldMarkFailedWithoutRetryWhenGRPCReturnsPermanentResult(t *testing.T) {
	repo := mocks.NewRepository(t)
	repo.EXPECT().ExistsProcessedEvent(mock.Anything, testdata.EventID).Return(false, nil)
	repo.EXPECT().CreateMintRequest(mock.Anything, mock.Anything).Return(testdata.PendingMintRequest(), nil)
	repo.EXPECT().CreateProcessedEvent(mock.Anything, testdata.EventID, "cashback.approved").Return(nil)
	repo.EXPECT().MarkFailed(mock.Anything, testdata.PendingMintRequest().ID, "invalid_address", "bad address", (*time.Time)(nil)).Return(nil)

	result := testdata.PermanentFailMintResult()
	bc := mocks.NewBlockchainClient(t)
	bc.EXPECT().MintToken(mock.Anything, mock.Anything).Return(result, nil)

	uc := mintcashback.NewUseCase(repo, bc, txManagerPassthrough(t))
	err := uc.Execute(context.Background(), testdata.NewInput())

	require.NoError(t, err)
}

func TestExecute_ShouldReturnErrorWhenExistsProcessedEventFails(t *testing.T) {
	repo := mocks.NewRepository(t)
	repo.EXPECT().ExistsProcessedEvent(mock.Anything, testdata.EventID).Return(false, errors.New("db error"))

	bc := mocks.NewBlockchainClient(t)

	uc := mintcashback.NewUseCase(repo, bc, txManagerPassthrough(t))
	err := uc.Execute(context.Background(), testdata.NewInput())

	require.Error(t, err)
}

func TestExecute_ShouldReturnErrorWhenCreateMintRequestFails(t *testing.T) {
	repo := mocks.NewRepository(t)
	repo.EXPECT().ExistsProcessedEvent(mock.Anything, testdata.EventID).Return(false, nil)
	repo.EXPECT().CreateMintRequest(mock.Anything, mock.Anything).Return(testdata.PendingMintRequest(), errors.New("constraint violation"))

	bc := mocks.NewBlockchainClient(t)

	uc := mintcashback.NewUseCase(repo, bc, txManagerPassthrough(t))
	err := uc.Execute(context.Background(), testdata.NewInput())

	require.Error(t, err)
}

func TestExecute_ShouldReturnErrorWhenCreateProcessedEventFails(t *testing.T) {
	repo := mocks.NewRepository(t)
	repo.EXPECT().ExistsProcessedEvent(mock.Anything, testdata.EventID).Return(false, nil)
	repo.EXPECT().CreateMintRequest(mock.Anything, mock.Anything).Return(testdata.PendingMintRequest(), nil)
	repo.EXPECT().CreateProcessedEvent(mock.Anything, testdata.EventID, "cashback.approved").Return(errors.New("unique constraint"))

	bc := mocks.NewBlockchainClient(t)

	uc := mintcashback.NewUseCase(repo, bc, txManagerPassthrough(t))
	err := uc.Execute(context.Background(), testdata.NewInput())

	require.Error(t, err)
}
