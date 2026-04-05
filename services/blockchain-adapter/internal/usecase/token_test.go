package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/cashback-platform/services/blockchain-adapter/internal/config"
	"github.com/cashback-platform/services/blockchain-adapter/internal/domain"
	"github.com/cashback-platform/services/blockchain-adapter/internal/usecase"
	"github.com/cashback-platform/services/blockchain-adapter/internal/usecase/mocks"
	ethereumpkg "github.com/cashback-platform/services/blockchain-adapter/pkg/ethereum"
)

const (
	testMnemonic       = "test test test test test test test test test test test junk"
	testDerivationPath = "m/44'/60'/0'/0/0"
	testWalletAddress  = "0x71C7656EC7ab88b098defB751B7401B5f6d8976F"
	testTokenAmount    = "1000000000000000000"
	testTxHash         = "0xabc123"
)

func newTestWallet(t *testing.T) *ethereumpkg.Wallet {
	t.Helper()
	w, err := ethereumpkg.NewFromMnemonic(testMnemonic, testDerivationPath)
	require.NoError(t, err)
	return w
}

func newTestConfig() *config.Config {
	return &config.Config{
		Ethereum: config.EthereumConfig{
			ChainID:         11155111,
			ContractAddress: "0x0000000000000000000000000000000000000000",
		},
	}
}

func TestMintToken_ShouldReturnExistingResultWhenAlreadySubmitted(t *testing.T) {
	nonceRepo := mocks.NewNonceRepository(t)
	txRepo := mocks.NewTransactionRepository(t)
	ethClient := mocks.NewEthereumClient(t)

	key := uuid.New()
	existing := &domain.BlockchainTransaction{
		ID:              1,
		IdempotencyKey:  key,
		Status:          domain.TransactionStatusSubmitted,
		TransactionHash: testTxHash,
	}
	txRepo.EXPECT().FindByIdempotencyKey(mock.Anything, key).Return(existing, nil)

	uc := usecase.NewToken(nonceRepo, txRepo, ethClient, newTestWallet(t), nil, newTestConfig())
	result, err := uc.MintToken(context.Background(), key.String(), testWalletAddress, testTokenAmount)

	require.NoError(t, err)
	require.True(t, result.Success)
	require.Equal(t, testTxHash, result.TransactionHash)
}

func TestMintToken_ShouldReturnErrorWhenTransactionAlreadyFailed(t *testing.T) {
	nonceRepo := mocks.NewNonceRepository(t)
	txRepo := mocks.NewTransactionRepository(t)
	ethClient := mocks.NewEthereumClient(t)

	key := uuid.New()
	existing := &domain.BlockchainTransaction{
		ID:             1,
		IdempotencyKey: key,
		Status:         domain.TransactionStatusFailed,
		ErrorCode:      "send_failed",
		ErrorMessage:   "timeout",
	}
	txRepo.EXPECT().FindByIdempotencyKey(mock.Anything, key).Return(existing, nil)

	uc := usecase.NewToken(nonceRepo, txRepo, ethClient, newTestWallet(t), nil, newTestConfig())
	result, err := uc.MintToken(context.Background(), key.String(), testWalletAddress, testTokenAmount)

	require.ErrorIs(t, err, usecase.ErrTransactionFailed)
	require.False(t, result.Success)
	require.False(t, result.Retryable)
}

func TestMintToken_ShouldReturnRetryableWhenPendingIsRecent(t *testing.T) {
	nonceRepo := mocks.NewNonceRepository(t)
	txRepo := mocks.NewTransactionRepository(t)
	ethClient := mocks.NewEthereumClient(t)

	key := uuid.New()
	existing := &domain.BlockchainTransaction{
		ID:             1,
		IdempotencyKey: key,
		Status:         domain.TransactionStatusPending,
		CreatedAt:      time.Now(),
	}
	txRepo.EXPECT().FindByIdempotencyKey(mock.Anything, key).Return(existing, nil)

	uc := usecase.NewToken(nonceRepo, txRepo, ethClient, newTestWallet(t), nil, newTestConfig())
	result, err := uc.MintToken(context.Background(), key.String(), testWalletAddress, testTokenAmount)

	require.NoError(t, err)
	require.False(t, result.Success)
	require.True(t, result.Retryable)
}

func TestMintToken_ShouldReturnLockUnavailableWhenNonceLockNotAcquired(t *testing.T) {
	nonceRepo := mocks.NewNonceRepository(t)
	txRepo := mocks.NewTransactionRepository(t)
	ethClient := mocks.NewEthereumClient(t)

	key := uuid.New()
	txRepo.EXPECT().FindByIdempotencyKey(mock.Anything, key).Return(nil, nil)
	nonceRepo.EXPECT().Increment(mock.Anything, testWalletAddress).Return(int64(0), errors.New("lock held"))

	uc := usecase.NewToken(nonceRepo, txRepo, ethClient, newTestWallet(t), nil, newTestConfig())
	result, err := uc.MintToken(context.Background(), key.String(), testWalletAddress, testTokenAmount)

	require.ErrorIs(t, err, usecase.ErrLockUnavailable)
	require.True(t, result.Retryable)
}

func TestMintToken_ShouldSyncNonceAndMarkFailedWhenSendFails(t *testing.T) {
	t.Skip("requires injecting pre-built signed tx; covered by integration test")
}

func TestMintToken_ShouldReturnSubmittedOnSuccess(t *testing.T) {
	// This test requires a real CashbackToken contract binding to call Mint() with NoSend=true.
	// Without a deployed contract or a stub, the abigen Mint() cannot produce a signed tx.
	// Covered by integration tests against a local Hardhat node.
	t.Skip("requires contract binding; covered by integration test")
}
