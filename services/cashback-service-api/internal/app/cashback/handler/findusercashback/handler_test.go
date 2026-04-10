package findusercashback_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	findusercashbackhandler "github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/handler/findusercashback"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/handler/findusercashback/testdata"
	findusercashbackuc "github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/usecase/findusercashback"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/usecase/findusercashback/mocks"
	"github.com/cashback-platform/kit/testsuite/handler"
)

type FindUserCashbackSuite struct {
	handler.Suite
}

func TestFindUserCashback(t *testing.T) {
	suite.Run(t, &FindUserCashbackSuite{})
}

func (s *FindUserCashbackSuite) TestSuccess() {
	s.Run("should return cashback summary", func() {
		t := s.T()

		repo := mocks.NewRepository(t)
		summary := testdata.CashbackSummary()
		repo.EXPECT().FindByUserID(mock.Anything, testdata.ValidUserID).Return(summary.Cashbacks, nil)
		repo.EXPECT().TotalByUserID(mock.Anything, testdata.ValidUserID).Return(summary.TotalMinted, nil)

		h := findusercashbackhandler.NewHandler(findusercashbackuc.New(repo))
		s.PrepareRouter(http.MethodGet, findusercashbackhandler.Path, h.Handle)
		s.Serve(findusercashbackhandler.Path, handler.WithPathParam("user_id", testdata.ValidUserID))

		resp := s.Response()
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body, `"total_cashbacks":2`)
	})
}

func (s *FindUserCashbackSuite) TestError() {
	s.Run("should return bad request when id is invalid", func() {
		t := s.T()

		repo := mocks.NewRepository(t)
		h := findusercashbackhandler.NewHandler(findusercashbackuc.New(repo))
		s.PrepareRouter(http.MethodGet, findusercashbackhandler.Path, h.Handle)
		s.Serve(findusercashbackhandler.Path, handler.WithPathParam("user_id", "abc"))

		resp := s.Response()
		require.Equal(t, http.StatusBadRequest, resp.Code)
	})

	s.Run("should return internal server error when use case fails", func() {
		t := s.T()

		repo := mocks.NewRepository(t)
		repo.EXPECT().FindByUserID(mock.Anything, testdata.ValidUserID).Return(nil, errors.New("db error"))

		h := findusercashbackhandler.NewHandler(findusercashbackuc.New(repo))
		s.PrepareRouter(http.MethodGet, findusercashbackhandler.Path, h.Handle)
		s.Serve(findusercashbackhandler.Path, handler.WithPathParam("user_id", testdata.ValidUserID))

		resp := s.Response()
		require.Equal(t, http.StatusInternalServerError, resp.Code)
	})
}
