package retrymints_test

import (
	"context"
	"errors"
	"testing"

	"github.com/cashback-platform/services/mint-consumer/internal/app/mint/domain"
	"github.com/cashback-platform/services/mint-consumer/internal/app/mint/usecase/retrymints"
	"github.com/cashback-platform/services/mint-consumer/internal/app/mint/usecase/retrymints/mocks"
	"github.com/cashback-platform/services/mint-consumer/internal/app/mint/usecase/retrymints/testdata"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestExecute_ShouldMarkCompletedWhenRetrySucceeds(t *testing.T) {
	repo := mocks.NewRepository(t)
	bc := mocks.NewBlockchainClient(t)

	result := testdata.SuccessfulMintResult()

	repo.EXPECT().FindFailedRetryable(mock.Anything, 50).Return([]domain.MintRequest{testdata.FailedRetryableMintRequest()}, nil)
	bc.EXPECT().MintToken(mock.Anything, mock.Anything).Return(result, nil)
	repo.EXPECT().MarkCompleted(mock.Anything, testdata.MintRequestID, result.TransactionHash, result.BlockNumber).Return(nil)

	uc := retrymints.NewUseCase(repo, bc)
	err := uc.Execute(context.Background())
	require.NoError(t, err)
}

func TestExecute_ShouldReturnNilWhenNoRetryableRequestsExist(t *testing.T) {
	repo := mocks.NewRepository(t)
	bc := mocks.NewBlockchainClient(t)

	repo.EXPECT().FindFailedRetryable(mock.Anything, 50).Return([]domain.MintRequest{}, nil)

	uc := retrymints.NewUseCase(repo, bc)
	err := uc.Execute(context.Background())
	require.NoError(t, err)
}

func TestExecute_ShouldMarkFailedWithNextRetryWhenGRPCReturnsRetryableResult(t *testing.T) {
	repo := mocks.NewRepository(t)
	bc := mocks.NewBlockchainClient(t)

	result := testdata.RetryableMintResult()

	repo.EXPECT().FindFailedRetryable(mock.Anything, 50).Return([]domain.MintRequest{testdata.FailedRetryableMintRequest()}, nil)
	bc.EXPECT().MintToken(mock.Anything, mock.Anything).Return(result, nil)
	repo.EXPECT().MarkFailed(mock.Anything, testdata.MintRequestID, result.ErrorCode, result.ErrorMessage, mock.Anything).Return(nil)

	uc := retrymints.NewUseCase(repo, bc)
	err := uc.Execute(context.Background())
	require.NoError(t, err)
}

func TestExecute_ShouldMarkFailedWithoutRetryWhenGRPCReturnsPermanentResult(t *testing.T) {
	repo := mocks.NewRepository(t)
	bc := mocks.NewBlockchainClient(t)

	result := testdata.PermanentFailMintResult()

	repo.EXPECT().FindFailedRetryable(mock.Anything, 50).Return([]domain.MintRequest{testdata.FailedRetryableMintRequest()}, nil)
	bc.EXPECT().MintToken(mock.Anything, mock.Anything).Return(result, nil)
	repo.EXPECT().MarkFailed(mock.Anything, testdata.MintRequestID, result.ErrorCode, result.ErrorMessage, mock.Anything).Return(nil)

	uc := retrymints.NewUseCase(repo, bc)
	err := uc.Execute(context.Background())
	require.NoError(t, err)
}

func TestExecute_ShouldReturnErrorWhenFindFailedRetryableFails(t *testing.T) {
	repo := mocks.NewRepository(t)
	bc := mocks.NewBlockchainClient(t)

	repo.EXPECT().FindFailedRetryable(mock.Anything, 50).Return(nil, errors.New("db error"))

	uc := retrymints.NewUseCase(repo, bc)
	err := uc.Execute(context.Background())
	require.Error(t, err)
}
