package balance_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/cashback-platform/kit/testsuite/handler"
	userdomain "github.com/cashback-platform/services/cashback-service-api/internal/app/user/domain"
	balancehandler "github.com/cashback-platform/services/cashback-service-api/internal/app/user/handler/balance"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/user/handler/balance/testdata"
	balanceuc "github.com/cashback-platform/services/cashback-service-api/internal/app/user/usecase/balance"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/user/usecase/balance/mocks"
)

type BalanceSuite struct {
	handler.Suite
}

func TestBalance(t *testing.T) {
	suite.Run(t, &BalanceSuite{})
}

func (s *BalanceSuite) TestSuccess() {
	s.Run("should return balance for valid user", func() {
		t := s.T()

		userRepo := mocks.NewUserRepository(t)
		blockchainClient := mocks.NewBlockchainClient(t)

		userRepo.EXPECT().FindByID(mock.Anything, testdata.ValidUserID).Return(testdata.FoundUser(), nil)
		blockchainClient.EXPECT().Balance(mock.Anything, testdata.WalletAddress).Return(testdata.TokenBalance(), nil)

		uc := balanceuc.New(userRepo, blockchainClient)
		h := balancehandler.NewHandler(uc)

		s.PrepareRouter(http.MethodGet, balancehandler.Path, h.Handle)
		s.Serve(balancehandler.Path, handler.WithPathParam("id", testdata.ValidUserID))

		resp := s.Response()
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body, `"balance_tokens":"1.00"`)
		require.Contains(t, resp.Body, `"wallet_address":"0xABCDEF1234567890"`)
	})
}

func (s *BalanceSuite) TestError() {
	s.Run("should return not found when user does not exist", func() {
		t := s.T()

		userRepo := mocks.NewUserRepository(t)
		blockchainClient := mocks.NewBlockchainClient(t)

		userRepo.EXPECT().FindByID(mock.Anything, testdata.ValidUserID).Return(userdomain.User{}, userdomain.ErrUserNotFound)

		uc := balanceuc.New(userRepo, blockchainClient)
		h := balancehandler.NewHandler(uc)

		s.PrepareRouter(http.MethodGet, balancehandler.Path, h.Handle)
		s.Serve(balancehandler.Path, handler.WithPathParam("id", testdata.ValidUserID))

		resp := s.Response()
		require.Equal(t, http.StatusNotFound, resp.Code)
	})

	s.Run("should return bad request when user id is not a number", func() {
		t := s.T()

		userRepo := mocks.NewUserRepository(t)
		blockchainClient := mocks.NewBlockchainClient(t)

		uc := balanceuc.New(userRepo, blockchainClient)
		h := balancehandler.NewHandler(uc)

		s.PrepareRouter(http.MethodGet, balancehandler.Path, h.Handle)
		s.Serve(balancehandler.Path, handler.WithPathParam("id", "abc"))

		resp := s.Response()
		require.Equal(t, http.StatusBadRequest, resp.Code)
	})

	s.Run("should return internal server error when blockchain fails", func() {
		t := s.T()

		userRepo := mocks.NewUserRepository(t)
		blockchainClient := mocks.NewBlockchainClient(t)

		userRepo.EXPECT().FindByID(mock.Anything, testdata.ValidUserID).Return(testdata.FoundUser(), nil)
		blockchainClient.EXPECT().Balance(mock.Anything, testdata.WalletAddress).Return(balanceuc.TokenBalance{}, errors.New("rpc error"))

		uc := balanceuc.New(userRepo, blockchainClient)
		h := balancehandler.NewHandler(uc)

		s.PrepareRouter(http.MethodGet, balancehandler.Path, h.Handle)
		s.Serve(balancehandler.Path, handler.WithPathParam("id", testdata.ValidUserID))

		resp := s.Response()
		require.Equal(t, http.StatusInternalServerError, resp.Code)
	})
}
