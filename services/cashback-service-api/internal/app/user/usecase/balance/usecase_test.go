package balance_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	userdomain "github.com/cashback-platform/services/cashback-service-api/internal/app/user/domain"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/user/usecase/balance"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/user/usecase/balance/mocks"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/user/usecase/balance/testdata"
)

func TestGetBalance_ShouldReturnBalanceWhenUserExists(t *testing.T) {
	userRepo := mocks.NewUserRepository(t)
	blockchainClient := mocks.NewBlockchainClient(t)

	userRepo.EXPECT().FindByID(mock.Anything, testdata.UserID).Return(testdata.FoundUser(), nil)
	blockchainClient.EXPECT().Balance(mock.Anything, testdata.WalletAddress).Return(testdata.TokenBalance(), nil)

	uc := balance.New(userRepo, blockchainClient)
	output, err := uc.Execute(t.Context(), testdata.UserID)

	require.NoError(t, err)
	require.Equal(t, testdata.UserID, output.UserID)
	require.Equal(t, testdata.WalletAddress, output.WalletAddress)
	require.Equal(t, "1000000000000000000", output.Balance)
	require.Equal(t, "1.00", output.BalanceTokens)
	require.Equal(t, int64(100), output.BlockNumber)
}

func TestGetBalance_ShouldReturnErrorWhenUserNotFound(t *testing.T) {
	userRepo := mocks.NewUserRepository(t)
	blockchainClient := mocks.NewBlockchainClient(t)

	userRepo.EXPECT().FindByID(mock.Anything, testdata.UserID).Return(userdomain.User{}, userdomain.ErrUserNotFound)

	uc := balance.New(userRepo, blockchainClient)
	_, err := uc.Execute(t.Context(), testdata.UserID)

	require.Error(t, err)
	require.ErrorIs(t, err, userdomain.ErrUserNotFound)
}

func TestGetBalance_ShouldReturnErrorWhenBlockchainFails(t *testing.T) {
	userRepo := mocks.NewUserRepository(t)
	blockchainClient := mocks.NewBlockchainClient(t)

	userRepo.EXPECT().FindByID(mock.Anything, testdata.UserID).Return(testdata.FoundUser(), nil)
	blockchainClient.EXPECT().Balance(mock.Anything, testdata.WalletAddress).Return(balance.TokenBalance{}, errors.New("rpc error"))

	uc := balance.New(userRepo, blockchainClient)
	_, err := uc.Execute(t.Context(), testdata.UserID)

	require.Error(t, err)
}

func TestGetBalance_ShouldReturnZeroTokensWhenBalanceIsZero(t *testing.T) {
	userRepo := mocks.NewUserRepository(t)
	blockchainClient := mocks.NewBlockchainClient(t)

	userRepo.EXPECT().FindByID(mock.Anything, testdata.UserID).Return(testdata.FoundUser(), nil)
	blockchainClient.EXPECT().Balance(mock.Anything, testdata.WalletAddress).Return(testdata.ZeroTokenBalance(), nil)

	uc := balance.New(userRepo, blockchainClient)
	output, err := uc.Execute(t.Context(), testdata.UserID)

	require.NoError(t, err)
	require.Equal(t, "0", output.BalanceTokens)
}

func TestGetBalance_ShouldReturnZeroWhenWeiBadFormat(t *testing.T) {
	userRepo := mocks.NewUserRepository(t)
	blockchainClient := mocks.NewBlockchainClient(t)

	userRepo.EXPECT().FindByID(mock.Anything, testdata.UserID).Return(testdata.FoundUser(), nil)
	blockchainClient.EXPECT().Balance(mock.Anything, testdata.WalletAddress).Return(testdata.InvalidWeiTokenBalance(), nil)

	uc := balance.New(userRepo, blockchainClient)
	output, err := uc.Execute(t.Context(), testdata.UserID)

	require.NoError(t, err)
	require.Equal(t, "0", output.BalanceTokens)
}
