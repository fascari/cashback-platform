package createuser_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/cashback-platform/kit/apperror"
	userdomain "github.com/cashback-platform/services/cashback-service-api/internal/app/user/domain"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/user/usecase/createuser"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/user/usecase/createuser/mocks"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/user/usecase/createuser/testdata"
)

func TestCreateUser_ShouldCreateUserWhenValid(t *testing.T) {
	repo := mocks.NewRepository(t)

	repo.EXPECT().FindByEmail(mock.Anything, testdata.Email).
		Return(userdomain.User{}, apperror.New(userdomain.ErrCodeUserNotFound, "not found"))
	repo.EXPECT().FindByExternalID(mock.Anything, testdata.ExternalID).
		Return(userdomain.User{}, apperror.New(userdomain.ErrCodeUserNotFound, "not found"))
	repo.EXPECT().Create(mock.Anything, mock.Anything).
		Return(testdata.CreatedUser(), nil)

	uc := createuser.New(repo)
	result, err := uc.Execute(t.Context(), testdata.ExternalID, testdata.Email, testdata.WalletAddress)

	require.NoError(t, err)
	require.Equal(t, testdata.Email, result.Email)
}

func TestCreateUser_ShouldReturnErrorWhenEmailAlreadyExists(t *testing.T) {
	repo := mocks.NewRepository(t)

	repo.EXPECT().FindByEmail(mock.Anything, testdata.Email).
		Return(userdomain.User{ID: 1}, nil)

	uc := createuser.New(repo)
	_, err := uc.Execute(t.Context(), testdata.ExternalID, testdata.Email, testdata.WalletAddress)

	require.ErrorIs(t, err, userdomain.ErrUserAlreadyExists)
}

func TestCreateUser_ShouldReturnErrorWhenExternalIDAlreadyExists(t *testing.T) {
	repo := mocks.NewRepository(t)

	repo.EXPECT().FindByEmail(mock.Anything, testdata.Email).
		Return(userdomain.User{}, apperror.New(userdomain.ErrCodeUserNotFound, "not found"))
	repo.EXPECT().FindByExternalID(mock.Anything, testdata.ExternalID).
		Return(userdomain.User{ID: 1}, nil)

	uc := createuser.New(repo)
	_, err := uc.Execute(t.Context(), testdata.ExternalID, testdata.Email, testdata.WalletAddress)

	require.ErrorIs(t, err, userdomain.ErrUserAlreadyExists)
}

func TestCreateUser_ShouldReturnErrorWhenRepositoryFails(t *testing.T) {
	repo := mocks.NewRepository(t)

	repo.EXPECT().FindByEmail(mock.Anything, testdata.Email).
		Return(userdomain.User{}, apperror.New(userdomain.ErrCodeUserNotFound, "not found"))
	repo.EXPECT().FindByExternalID(mock.Anything, testdata.ExternalID).
		Return(userdomain.User{}, apperror.New(userdomain.ErrCodeUserNotFound, "not found"))
	repo.EXPECT().Create(mock.Anything, mock.Anything).
		Return(userdomain.User{}, errors.New("db error"))

	uc := createuser.New(repo)
	_, err := uc.Execute(t.Context(), testdata.ExternalID, testdata.Email, testdata.WalletAddress)

	require.Error(t, err)
}
