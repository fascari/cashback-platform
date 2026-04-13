package calculatecashback_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/cashback-platform/kit/apperror"
	cashdomain "github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/domain"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/usecase/calculatecashback"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/usecase/calculatecashback/mocks"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/usecase/calculatecashback/testdata"
	purchasedomain "github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/domain"
	userdomain "github.com/cashback-platform/services/cashback-service-api/internal/app/user/domain"
)

func TestCalculateCashback_ShouldCalculateAndSaveCashback(t *testing.T) {
	repo := mocks.NewRepository(t)
	purchaseRepo := mocks.NewPurchaseRepository(t)
	userRepo := mocks.NewUserRepository(t)

	repo.EXPECT().FindByPurchaseID(mock.Anything, testdata.PurchaseID).
		Return(cashdomain.Cashback{}, apperror.New(cashdomain.ErrCodeCashbackNotFound, "not found"))
	purchaseRepo.EXPECT().FindByID(mock.Anything, testdata.PurchaseID).
		Return(testdata.NewPurchase(), nil)
	userRepo.EXPECT().FindByID(mock.Anything, testdata.UserID).
		Return(testdata.NewUser(), nil)
	repo.EXPECT().CreateWithEvent(mock.Anything, mock.Anything, mock.Anything).
		Return(testdata.ApprovedCashback(), nil)

	uc := calculatecashback.New(repo, purchaseRepo, userRepo)
	result, err := uc.Execute(t.Context(), testdata.PurchaseID)

	require.NoError(t, err)
	require.Equal(t, cashdomain.StatusApproved, result.Status)
}

func TestCalculateCashback_ShouldReturnErrorWhenCashbackAlreadyExists(t *testing.T) {
	repo := mocks.NewRepository(t)
	purchaseRepo := mocks.NewPurchaseRepository(t)
	userRepo := mocks.NewUserRepository(t)

	repo.EXPECT().FindByPurchaseID(mock.Anything, testdata.PurchaseID).
		Return(testdata.ExistingCashback(), nil)

	uc := calculatecashback.New(repo, purchaseRepo, userRepo)
	_, err := uc.Execute(t.Context(), testdata.PurchaseID)

	require.ErrorIs(t, err, cashdomain.ErrCashbackAlreadyExists)
}

func TestCalculateCashback_ShouldReturnErrorWhenPurchaseNotFound(t *testing.T) {
	repo := mocks.NewRepository(t)
	purchaseRepo := mocks.NewPurchaseRepository(t)
	userRepo := mocks.NewUserRepository(t)

	repo.EXPECT().FindByPurchaseID(mock.Anything, testdata.PurchaseID).
		Return(cashdomain.Cashback{}, apperror.New(cashdomain.ErrCodeCashbackNotFound, "not found"))
	purchaseRepo.EXPECT().FindByID(mock.Anything, testdata.PurchaseID).
		Return(purchasedomain.Purchase{}, errors.New("not found"))

	uc := calculatecashback.New(repo, purchaseRepo, userRepo)
	_, err := uc.Execute(t.Context(), testdata.PurchaseID)

	require.ErrorIs(t, err, calculatecashback.ErrPurchaseNotFound)
}

func TestCalculateCashback_ShouldReturnErrorWhenUserNotFound(t *testing.T) {
	repo := mocks.NewRepository(t)
	purchaseRepo := mocks.NewPurchaseRepository(t)
	userRepo := mocks.NewUserRepository(t)

	repo.EXPECT().FindByPurchaseID(mock.Anything, testdata.PurchaseID).
		Return(cashdomain.Cashback{}, apperror.New(cashdomain.ErrCodeCashbackNotFound, "not found"))
	purchaseRepo.EXPECT().FindByID(mock.Anything, testdata.PurchaseID).
		Return(testdata.NewPurchase(), nil)
	userRepo.EXPECT().FindByID(mock.Anything, testdata.UserID).
		Return(userdomain.User{}, errors.New("not found"))

	uc := calculatecashback.New(repo, purchaseRepo, userRepo)
	_, err := uc.Execute(t.Context(), testdata.PurchaseID)

	require.ErrorIs(t, err, calculatecashback.ErrUserNotFound)
}

func TestCalculateCashback_ShouldReturnErrorWhenCreateEventFails(t *testing.T) {
	repo := mocks.NewRepository(t)
	purchaseRepo := mocks.NewPurchaseRepository(t)
	userRepo := mocks.NewUserRepository(t)

	repo.EXPECT().FindByPurchaseID(mock.Anything, testdata.PurchaseID).
		Return(cashdomain.Cashback{}, apperror.New(cashdomain.ErrCodeCashbackNotFound, "not found"))
	purchaseRepo.EXPECT().FindByID(mock.Anything, testdata.PurchaseID).
		Return(testdata.NewPurchase(), nil)
	userRepo.EXPECT().FindByID(mock.Anything, testdata.UserID).
		Return(testdata.NewUser(), nil)
	repo.EXPECT().CreateWithEvent(mock.Anything, mock.Anything, mock.Anything).
		Return(cashdomain.Cashback{}, errors.New("db error"))

	uc := calculatecashback.New(repo, purchaseRepo, userRepo)
	_, err := uc.Execute(t.Context(), testdata.PurchaseID)

	require.Error(t, err)
}
