package usecase_test

import (
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/cashback-platform/services/blockchain-adapter/internal/app/token/domain"
	"github.com/cashback-platform/services/blockchain-adapter/internal/app/token/usecase"
	"github.com/cashback-platform/services/blockchain-adapter/internal/app/token/usecase/mocks"
	"github.com/cashback-platform/services/blockchain-adapter/internal/app/token/usecase/testdata"
)

func TestMintToken_ShouldReturnExistingResultWhenAlreadySubmitted(t *testing.T) {
	nonceRepo := mocks.NewNonceRepository(t)
	txRepo := mocks.NewTransactionRepository(t)
	ethClient := mocks.NewEthereumClient(t)
	token := mocks.NewTokenContract(t)

	key := uuid.New()
	txRepo.EXPECT().FindByIdempotencyKey(mock.Anything, key).Return(testdata.SubmittedTransaction(key), nil)

	uc := usecase.NewToken(nonceRepo, txRepo, ethClient, testdata.Wallet(t), token, testdata.Config())
	result, err := uc.MintToken(t.Context(), key.String(), testdata.WalletAddress, testdata.TokenAmount)

	require.NoError(t, err)
	require.True(t, result.Success)
	require.Equal(t, testdata.TxHash, result.TransactionHash)
}

func TestMintToken_ShouldReturnErrorWhenTransactionAlreadyFailed(t *testing.T) {
	nonceRepo := mocks.NewNonceRepository(t)
	txRepo := mocks.NewTransactionRepository(t)
	ethClient := mocks.NewEthereumClient(t)
	token := mocks.NewTokenContract(t)

	key := uuid.New()
	txRepo.EXPECT().FindByIdempotencyKey(mock.Anything, key).Return(testdata.FailedTransaction(key), nil)

	uc := usecase.NewToken(nonceRepo, txRepo, ethClient, testdata.Wallet(t), token, testdata.Config())
	result, err := uc.MintToken(t.Context(), key.String(), testdata.WalletAddress, testdata.TokenAmount)

	require.ErrorIs(t, err, usecase.ErrTransactionFailed)
	require.False(t, result.Success)
	require.False(t, result.Retryable)
}

func TestMintToken_ShouldReturnRetryableWhenPendingIsRecent(t *testing.T) {
	nonceRepo := mocks.NewNonceRepository(t)
	txRepo := mocks.NewTransactionRepository(t)
	ethClient := mocks.NewEthereumClient(t)
	token := mocks.NewTokenContract(t)

	key := uuid.New()
	txRepo.EXPECT().FindByIdempotencyKey(mock.Anything, key).Return(testdata.RecentPendingTransaction(key), nil)

	uc := usecase.NewToken(nonceRepo, txRepo, ethClient, testdata.Wallet(t), token, testdata.Config())
	result, err := uc.MintToken(t.Context(), key.String(), testdata.WalletAddress, testdata.TokenAmount)

	require.NoError(t, err)
	require.False(t, result.Success)
	require.True(t, result.Retryable)
}

func TestMintToken_ShouldReturnLockUnavailableWhenNonceLockNotAcquired(t *testing.T) {
	nonceRepo := mocks.NewNonceRepository(t)
	txRepo := mocks.NewTransactionRepository(t)
	ethClient := mocks.NewEthereumClient(t)
	token := mocks.NewTokenContract(t)

	key := uuid.New()
	txRepo.EXPECT().FindByIdempotencyKey(mock.Anything, key).Return(nil, nil)
	nonceRepo.EXPECT().Increment(mock.Anything, testdata.WalletAddress).Return(int64(0), errors.New("lock held"))

	uc := usecase.NewToken(nonceRepo, txRepo, ethClient, testdata.Wallet(t), token, testdata.Config())
	result, err := uc.MintToken(t.Context(), key.String(), testdata.WalletAddress, testdata.TokenAmount)

	require.ErrorIs(t, err, usecase.ErrLockUnavailable)
	require.True(t, result.Retryable)
}

func TestMintToken_ShouldSyncNonceAndMarkFailedWhenSendFails(t *testing.T) {
	nonceRepo := mocks.NewNonceRepository(t)
	txRepo := mocks.NewTransactionRepository(t)
	ethClient := mocks.NewEthereumClient(t)
	token := mocks.NewTokenContract(t)

	key := uuid.New()
	sendErr := errors.New("connection refused")
	tx := testdata.SignedTransaction()

	txRepo.EXPECT().FindByIdempotencyKey(mock.Anything, key).Return(nil, nil)
	nonceRepo.EXPECT().Increment(mock.Anything, testdata.WalletAddress).Return(testdata.Nonce, nil)
	ethClient.EXPECT().SuggestGasPrice(mock.Anything).Return(big.NewInt(1e9), nil)
	token.EXPECT().Mint(mock.Anything, common.HexToAddress(testdata.WalletAddress), mock.Anything).Return(tx, nil)
	txRepo.EXPECT().Create(mock.Anything, mock.Anything).Return(domain.BlockchainTransaction{ID: testdata.RecordID}, nil)
	ethClient.EXPECT().SendTransaction(mock.Anything, tx).Return(sendErr).Times(3)
	ethClient.EXPECT().PendingNonceAt(mock.Anything, mock.Anything).Return(uint64(testdata.OnChainNonce), nil)
	nonceRepo.EXPECT().SyncFromChain(mock.Anything, testdata.WalletAddress, testdata.OnChainNonce).Return(nil)
	txRepo.EXPECT().MarkFailed(mock.Anything, testdata.RecordID, "send_failed", sendErr.Error()).Return(nil)

	uc := usecase.NewToken(nonceRepo, txRepo, ethClient, testdata.Wallet(t), token, testdata.Config())
	result, err := uc.MintToken(t.Context(), key.String(), testdata.WalletAddress, testdata.TokenAmount)

	require.NoError(t, err)
	require.False(t, result.Success)
	require.True(t, result.Retryable)
	require.Equal(t, "send_failed", result.ErrorCode)
}

func TestMintToken_ShouldReturnSubmittedOnSuccess(t *testing.T) {
	nonceRepo := mocks.NewNonceRepository(t)
	txRepo := mocks.NewTransactionRepository(t)
	ethClient := mocks.NewEthereumClient(t)
	token := mocks.NewTokenContract(t)

	key := uuid.New()
	tx := testdata.SignedTransaction()

	txRepo.EXPECT().FindByIdempotencyKey(mock.Anything, key).Return(nil, nil)
	nonceRepo.EXPECT().Increment(mock.Anything, testdata.WalletAddress).Return(testdata.Nonce, nil)
	ethClient.EXPECT().SuggestGasPrice(mock.Anything).Return(big.NewInt(1e9), nil)
	token.EXPECT().Mint(mock.Anything, common.HexToAddress(testdata.WalletAddress), mock.Anything).Return(tx, nil)
	txRepo.EXPECT().Create(mock.Anything, mock.Anything).Return(domain.BlockchainTransaction{ID: testdata.RecordID}, nil)
	ethClient.EXPECT().SendTransaction(mock.Anything, tx).Return(nil)
	txRepo.EXPECT().UpdateStatus(mock.Anything, testdata.RecordID, domain.TransactionStatusSubmitted, tx.Hash().Hex(), int64(0)).Return(nil)

	uc := usecase.NewToken(nonceRepo, txRepo, ethClient, testdata.Wallet(t), token, testdata.Config())
	result, err := uc.MintToken(t.Context(), key.String(), testdata.WalletAddress, testdata.TokenAmount)

	require.NoError(t, err)
	require.True(t, result.Success)
	require.Equal(t, tx.Hash().Hex(), result.TransactionHash)
	require.Equal(t, string(domain.TransactionStatusSubmitted), result.Status)
}

