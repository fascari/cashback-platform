package finduser_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/cashback-platform/kit/testsuite/handler"
	userdomain "github.com/cashback-platform/services/cashback-service-api/internal/app/user/domain"
	finduserhandler "github.com/cashback-platform/services/cashback-service-api/internal/app/user/handler/finduser"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/user/handler/finduser/testdata"
	finduseruc "github.com/cashback-platform/services/cashback-service-api/internal/app/user/usecase/finduser"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/user/usecase/finduser/mocks"
)

type FindUserSuite struct {
	handler.Suite
}

func TestFindUser(t *testing.T) {
	suite.Run(t, &FindUserSuite{})
}

func (s *FindUserSuite) TestSuccess() {
	s.Run("should return user when found", func() {
		t := s.T()

		repo := mocks.NewRepository(t)
		repo.EXPECT().FindByID(mock.Anything, testdata.ValidID).Return(testdata.ExistingUser(), nil)

		h := finduserhandler.NewHandler(finduseruc.New(repo))
		s.PrepareRouter(http.MethodGet, finduserhandler.Path, h.Handle)
		s.Serve(finduserhandler.Path, handler.WithPathParam("id", testdata.ValidID))

		resp := s.Response()
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body, `"email":"user@example.com"`)
	})
}

func (s *FindUserSuite) TestError() {
	s.Run("should return not found when user does not exist", func() {
		t := s.T()

		repo := mocks.NewRepository(t)
		repo.EXPECT().FindByID(mock.Anything, testdata.ValidID).Return(userdomain.User{}, userdomain.ErrUserNotFound)

		h := finduserhandler.NewHandler(finduseruc.New(repo))
		s.PrepareRouter(http.MethodGet, finduserhandler.Path, h.Handle)
		s.Serve(finduserhandler.Path, handler.WithPathParam("id", testdata.ValidID))

		resp := s.Response()
		require.Equal(t, http.StatusNotFound, resp.Code)
	})

	s.Run("should return bad request when id is invalid", func() {
		t := s.T()

		repo := mocks.NewRepository(t)
		h := finduserhandler.NewHandler(finduseruc.New(repo))
		s.PrepareRouter(http.MethodGet, finduserhandler.Path, h.Handle)
		s.Serve(finduserhandler.Path, handler.WithPathParam("id", "abc"))

		resp := s.Response()
		require.Equal(t, http.StatusBadRequest, resp.Code)
	})

	s.Run("should return internal server error when use case fails", func() {
		t := s.T()

		repo := mocks.NewRepository(t)
		repo.EXPECT().FindByID(mock.Anything, testdata.ValidID).Return(userdomain.User{}, errors.New("db error"))

		h := finduserhandler.NewHandler(finduseruc.New(repo))
		s.PrepareRouter(http.MethodGet, finduserhandler.Path, h.Handle)
		s.Serve(finduserhandler.Path, handler.WithPathParam("id", testdata.ValidID))

		resp := s.Response()
		require.Equal(t, http.StatusInternalServerError, resp.Code)
	})
}
