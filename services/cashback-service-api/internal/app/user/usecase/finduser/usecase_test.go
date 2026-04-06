package finduser_test

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	userdomain "github.com/cashback-platform/services/cashback-service-api/internal/app/user/domain"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/user/usecase/finduser"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/user/usecase/finduser/mocks"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/user/usecase/finduser/testdata"
)

func TestFindUser_ShouldReturnUserWhenFound(t *testing.T) {
	repo := mocks.NewRepository(t)

	repo.EXPECT().FindByID(mock.Anything, testdata.UserID).Return(testdata.FoundUser(), nil)

	uc := finduser.New(repo)
	result, err := uc.Execute(t.Context(), testdata.UserID)

	require.NoError(t, err)
	require.Equal(t, testdata.FoundUser(), result)
}

func TestFindUser_ShouldReturnErrorWhenNotFound(t *testing.T) {
	repo := mocks.NewRepository(t)

	const id int64 = 99

	repo.EXPECT().FindByID(mock.Anything, id).
		Return(userdomain.User{}, userdomain.ErrUserNotFound)

	uc := finduser.New(repo)
	_, err := uc.Execute(t.Context(), id)

	require.ErrorIs(t, err, userdomain.ErrUserNotFound)
}
