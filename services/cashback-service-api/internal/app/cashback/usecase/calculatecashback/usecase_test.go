package calculatecashback_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	cashdomain "github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/domain"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/usecase/calculatecashback"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/usecase/calculatecashback/mocks"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/usecase/calculatecashback/testdata"
	"github.com/cashback-platform/services/cashback-service-api/pkg/apperror"
)

func TestCalculateCashback_ShouldCalculateAndSaveCashback(t *testing.T) {
	repo := mocks.NewRepository(t)
	purchaseRepo := mocks.NewPurchaseRepository(t)
	userRepo := mocks.NewUserRepository(t)
	publisher := mocks.NewEventPublisher(t)

	repo.EXPECT().FindByPurchaseID(mock.Anything, testdata.PurchaseID).
		Return(cashdomain.Cashback{}, apperror.New(cashdomain.ErrCodeCashbackNotFound, "not found"))
	purchaseRepo.EXPECT().FindByID(mock.Anything, testdata.PurchaseID).
		Return(testdata.NewPurchase(), nil)
	userRepo.EXPECT().FindByID(mock.Anything, testdata.UserID).
		Return(testdata.NewUser(), nil)
	repo.EXPECT().Create(mock.Anything, mock.Anything).
		Return(testdata.ApprovedCashback(), nil)
	publisher.EXPECT().Publish(mock.Anything, calculatecashback.EventTypeCashbackApproved, mock.Anything, mock.Anything).
		Return(nil)

	uc := calculatecashback.New(repo, purchaseRepo, userRepo, publisher)
	result, err := uc.Execute(context.Background(), testdata.PurchaseID)

	require.NoError(t, err)
	require.Equal(t, cashdomain.StatusApproved, result.Status)
}

func TestCalculateCashback_ShouldReturnErrorWhenCashbackAlreadyExists(t *testing.T) {
	repo := mocks.NewRepository(t)
	purchaseRepo := mocks.NewPurchaseRepository(t)
	userRepo := mocks.NewUserRepository(t)
	publisher := mocks.NewEventPublisher(t)

	repo.EXPECT().FindByPurchaseID(mock.Anything, testdata.PurchaseID).
		Return(testdata.ExistingCashback(), nil)

	uc := calculatecashback.New(repo, purchaseRepo, userRepo, publisher)
	_, err := uc.Execute(context.Background(), testdata.PurchaseID)

	require.ErrorIs(t, err, cashdomain.ErrCashbackAlreadyExists)
}

func TestCalculateCashback_ShouldReturnErrorWhenPurchaseNotFound(t *testing.T) {
	repo := mocks.NewRepository(t)
	purchaseRepo := mocks.NewPurchaseRepository(t)
	userRepo := mocks.NewUserRepository(t)
	publisher := mocks.NewEventPublisher(t)

	repo.EXPECT().FindByPurchaseID(mock.Anything, testdata.PurchaseID).
		Return(cashdomain.Cashback{}, apperror.New(cashdomain.ErrCodeCashbackNotFound, "not found"))
	purchaseRepo.EXPECT().FindByID(mock.Anything, testdata.PurchaseID).
		Return(calculatecashback.Purchase{}, errors.New("not found"))

	uc := calculatecashback.New(repo, purchaseRepo, userRepo, publisher)
	_, err := uc.Execute(context.Background(), testdata.PurchaseID)

	require.ErrorIs(t, err, calculatecashback.ErrPurchaseNotFound)
}

func TestCalculateCashback_ShouldReturnErrorWhenUserNotFound(t *testing.T) {
	repo := mocks.NewRepository(t)
	purchaseRepo := mocks.NewPurchaseRepository(t)
	userRepo := mocks.NewUserRepository(t)
	publisher := mocks.NewEventPublisher(t)

	repo.EXPECT().FindByPurchaseID(mock.Anything, testdata.PurchaseID).
		Return(cashdomain.Cashback{}, apperror.New(cashdomain.ErrCodeCashbackNotFound, "not found"))
	purchaseRepo.EXPECT().FindByID(mock.Anything, testdata.PurchaseID).
		Return(testdata.NewPurchase(), nil)
	userRepo.EXPECT().FindByID(mock.Anything, testdata.UserID).
		Return(calculatecashback.User{}, errors.New("not found"))

	uc := calculatecashback.New(repo, purchaseRepo, userRepo, publisher)
	_, err := uc.Execute(context.Background(), testdata.PurchaseID)

	require.ErrorIs(t, err, calculatecashback.ErrUserNotFound)
}

func TestCalculateCashback_ShouldReturnErrorWhenPublishFails(t *testing.T) {
	repo := mocks.NewRepository(t)
	purchaseRepo := mocks.NewPurchaseRepository(t)
	userRepo := mocks.NewUserRepository(t)
	publisher := mocks.NewEventPublisher(t)

	repo.EXPECT().FindByPurchaseID(mock.Anything, testdata.PurchaseID).
		Return(cashdomain.Cashback{}, apperror.New(cashdomain.ErrCodeCashbackNotFound, "not found"))
	purchaseRepo.EXPECT().FindByID(mock.Anything, testdata.PurchaseID).
		Return(testdata.NewPurchase(), nil)
	userRepo.EXPECT().FindByID(mock.Anything, testdata.UserID).
		Return(testdata.NewUser(), nil)
	repo.EXPECT().Create(mock.Anything, mock.Anything).
		Return(testdata.ApprovedCashback(), nil)
	publisher.EXPECT().Publish(mock.Anything, calculatecashback.EventTypeCashbackApproved, mock.Anything, mock.Anything).
		Return(errors.New("broker unavailable"))

	uc := calculatecashback.New(repo, purchaseRepo, userRepo, publisher)
	_, err := uc.Execute(context.Background(), testdata.PurchaseID)

	require.ErrorIs(t, err, calculatecashback.ErrFailedToPublishEvent)
}
