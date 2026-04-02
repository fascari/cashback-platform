package createuser_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	userdomain "github.com/cashback-platform/services/cashback-service-api/internal/app/user/domain"
	createuserhandler "github.com/cashback-platform/services/cashback-service-api/internal/app/user/handler/createuser"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/user/handler/createuser/testdata"
	createuseruc "github.com/cashback-platform/services/cashback-service-api/internal/app/user/usecase/createuser"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/user/usecase/createuser/mocks"
	"github.com/cashback-platform/services/cashback-service-api/pkg/testsuite/handler"
)

type CreateUserSuite struct {
	handler.Suite
}

func TestCreateUser(t *testing.T) {
	suite.Run(t, &CreateUserSuite{})
}

func (s *CreateUserSuite) TestSuccess() {
	s.Run("should return created user", func() {
		t := s.T()

		repo := mocks.NewRepository(t)
		repo.EXPECT().FindByEmail(mock.Anything, testdata.ValidPayload().Email).Return(userdomain.User{}, userdomain.ErrUserNotFound)
		repo.EXPECT().FindByExternalID(mock.Anything, testdata.ValidPayload().ExternalID).Return(userdomain.User{}, userdomain.ErrUserNotFound)
		repo.EXPECT().Create(mock.Anything, mock.Anything).Return(testdata.CreatedUser(), nil)

		h := createuserhandler.NewHandler(createuseruc.New(repo))
		s.PrepareRouter(http.MethodPost, createuserhandler.Path, h.Handle)
		s.Serve(createuserhandler.Path, handler.WithJSONBodyStruct(testdata.ValidPayload()))

		resp := s.Response()
		require.Equal(t, http.StatusCreated, resp.Code)
		require.Contains(t, resp.Body, `"email":"user@example.com"`)
	})
}

func (s *CreateUserSuite) TestError() {
	s.Run("should return bad request when email is invalid", func() {
		t := s.T()

		repo := mocks.NewRepository(t)
		h := createuserhandler.NewHandler(createuseruc.New(repo))
		s.PrepareRouter(http.MethodPost, createuserhandler.Path, h.Handle)
		s.Serve(createuserhandler.Path, handler.WithJSONBodyStruct(testdata.InvalidEmailPayload()))

		resp := s.Response()
		require.Equal(t, http.StatusBadRequest, resp.Code)
		require.Contains(t, resp.Body, "email must be a valid email address")
	})

	s.Run("should return bad request when wallet address is too short", func() {
		t := s.T()

		repo := mocks.NewRepository(t)
		h := createuserhandler.NewHandler(createuseruc.New(repo))
		s.PrepareRouter(http.MethodPost, createuserhandler.Path, h.Handle)
		s.Serve(createuserhandler.Path, handler.WithJSONBodyStruct(testdata.ShortWalletPayload()))

		resp := s.Response()
		require.Equal(t, http.StatusBadRequest, resp.Code)
		require.Contains(t, resp.Body, "wallet_address must be at least 20 characters in length")
	})

	s.Run("should return conflict when user already exists", func() {
		t := s.T()

		repo := mocks.NewRepository(t)
		repo.EXPECT().FindByEmail(mock.Anything, testdata.ValidPayload().Email).Return(testdata.ExistingUser(), nil)

		h := createuserhandler.NewHandler(createuseruc.New(repo))
		s.PrepareRouter(http.MethodPost, createuserhandler.Path, h.Handle)
		s.Serve(createuserhandler.Path, handler.WithJSONBodyStruct(testdata.ValidPayload()))

		resp := s.Response()
		require.Equal(t, http.StatusConflict, resp.Code)
	})

	s.Run("should return internal server error when use case fails", func() {
		t := s.T()

		repo := mocks.NewRepository(t)
		repo.EXPECT().FindByEmail(mock.Anything, testdata.ValidPayload().Email).Return(userdomain.User{}, errors.New("db error"))

		h := createuserhandler.NewHandler(createuseruc.New(repo))
		s.PrepareRouter(http.MethodPost, createuserhandler.Path, h.Handle)
		s.Serve(createuserhandler.Path, handler.WithJSONBodyStruct(testdata.ValidPayload()))

		resp := s.Response()
		require.Equal(t, http.StatusInternalServerError, resp.Code)
	})

	s.Run("should return bad request when body is invalid JSON", func() {
		t := s.T()

		repo := mocks.NewRepository(t)
		h := createuserhandler.NewHandler(createuseruc.New(repo))
		s.PrepareRouter(http.MethodPost, createuserhandler.Path, h.Handle)
		s.Serve(createuserhandler.Path, handler.WithJSONBody("not valid json"))

		resp := s.Response()
		require.Equal(t, http.StatusBadRequest, resp.Code)
	})
}
