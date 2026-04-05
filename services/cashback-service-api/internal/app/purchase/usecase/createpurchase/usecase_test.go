package createpurchase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	purchasedomain "github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/domain"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/usecase/createpurchase"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/usecase/createpurchase/mocks"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/usecase/createpurchase/testdata"
)

func TestCreatePurchase_ShouldCreatePurchaseWhenValid(t *testing.T) {
	repo := mocks.NewRepository(t)

	repo.EXPECT().Create(mock.Anything, mock.Anything).Return(testdata.CreatedPurchase(), nil)

	uc := createpurchase.New(repo)
	result, err := uc.Execute(context.Background(), testdata.UserID, testdata.Amount, testdata.MerchantID)

	require.NoError(t, err)
	require.Equal(t, testdata.CreatedPurchase().ID, result.ID)
}

func TestCreatePurchase_ShouldReturnErrorWhenAmountIsZero(t *testing.T) {
	repo := mocks.NewRepository(t)

	uc := createpurchase.New(repo)
	_, err := uc.Execute(context.Background(), 1, 0, "merchant")

	require.ErrorIs(t, err, createpurchase.ErrInvalidAmount)
}

func TestCreatePurchase_ShouldReturnErrorWhenAmountIsNegative(t *testing.T) {
	repo := mocks.NewRepository(t)

	uc := createpurchase.New(repo)
	_, err := uc.Execute(context.Background(), 1, -10.0, "merchant")

	require.ErrorIs(t, err, createpurchase.ErrInvalidAmount)
}

func TestCreatePurchase_ShouldReturnErrorWhenUserIDIsZero(t *testing.T) {
	repo := mocks.NewRepository(t)

	uc := createpurchase.New(repo)
	_, err := uc.Execute(context.Background(), 0, 100.0, "merchant")

	require.ErrorIs(t, err, createpurchase.ErrInvalidUserID)
}

func TestCreatePurchase_ShouldReturnErrorWhenMerchantIsEmpty(t *testing.T) {
	repo := mocks.NewRepository(t)

	uc := createpurchase.New(repo)
	_, err := uc.Execute(context.Background(), 1, 100.0, "")

	require.ErrorIs(t, err, createpurchase.ErrInvalidMerchant)
}

func TestCreatePurchase_ShouldReturnErrorWhenRepositoryFails(t *testing.T) {
	repo := mocks.NewRepository(t)

	repo.EXPECT().Create(mock.Anything, mock.Anything).Return(purchasedomain.Purchase{}, errors.New("db error"))

	uc := createpurchase.New(repo)
	_, err := uc.Execute(context.Background(), 1, 100.0, "merchant")

	require.Error(t, err)
}
