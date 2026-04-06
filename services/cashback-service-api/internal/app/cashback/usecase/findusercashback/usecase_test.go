package findusercashback_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	cashdomain "github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/domain"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/usecase/findusercashback"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/usecase/findusercashback/mocks"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/usecase/findusercashback/testdata"
)

func TestFindUserCashback_ShouldReturnCashbacksForUser(t *testing.T) {
	repo := mocks.NewRepository(t)

	repo.EXPECT().FindByUserID(mock.Anything, testdata.UserID).Return(testdata.Cashbacks(), nil)
	repo.EXPECT().TotalByUserID(mock.Anything, testdata.UserID).Return(8.0, nil)

	uc := findusercashback.New(repo)
	result, err := uc.Execute(t.Context(), testdata.UserID)

	require.NoError(t, err)
	require.Len(t, result.Cashbacks, 2)
	require.Equal(t, 8.0, result.TotalMinted)
	require.Equal(t, 2, result.TotalCashbacks)
}

func TestFindUserCashback_ShouldReturnEmptyListWhenNoCashbacks(t *testing.T) {
	repo := mocks.NewRepository(t)

	repo.EXPECT().FindByUserID(mock.Anything, testdata.UserID).Return([]cashdomain.Cashback{}, nil)
	repo.EXPECT().TotalByUserID(mock.Anything, testdata.UserID).Return(0.0, nil)

	uc := findusercashback.New(repo)
	result, err := uc.Execute(t.Context(), testdata.UserID)

	require.NoError(t, err)
	require.Empty(t, result.Cashbacks)
	require.Equal(t, 0, result.TotalCashbacks)
}

func TestFindUserCashback_ShouldReturnErrorWhenUserIDIsZero(t *testing.T) {
	repo := mocks.NewRepository(t)

	uc := findusercashback.New(repo)
	_, err := uc.Execute(t.Context(), 0)

	require.ErrorIs(t, err, cashdomain.ErrInvalidUserID)
}

func TestFindUserCashback_ShouldReturnErrorWhenRepositoryFails(t *testing.T) {
	repo := mocks.NewRepository(t)

	repo.EXPECT().FindByUserID(mock.Anything, testdata.UserID).Return(nil, errors.New("db error"))

	uc := findusercashback.New(repo)
	_, err := uc.Execute(t.Context(), testdata.UserID)

	require.Error(t, err)
}
