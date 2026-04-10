package findpurchase_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/cashback-platform/kit/testsuite/handler"
	purchasedomain "github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/domain"
	findpurchasehandler "github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/handler/findpurchase"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/handler/findpurchase/testdata"
	findpurchaseuc "github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/usecase/findpurchase"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/usecase/findpurchase/mocks"
)

type FindPurchaseSuite struct {
	handler.Suite
}

func TestFindPurchase(t *testing.T) {
	suite.Run(t, &FindPurchaseSuite{})
}

func (s *FindPurchaseSuite) TestSuccess() {
	s.Run("should return purchase when found", func() {
		t := s.T()

		repo := mocks.NewRepository(t)
		repo.EXPECT().FindByID(mock.Anything, testdata.ValidID).Return(testdata.ExistingPurchase(), nil)

		h := findpurchasehandler.NewHandler(findpurchaseuc.New(repo))
		s.PrepareRouter(http.MethodGet, findpurchasehandler.Path, h.Handle)
		s.Serve(findpurchasehandler.Path, handler.WithPathParam("id", testdata.ValidID))

		resp := s.Response()
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body, `"status":"pending"`)
	})
}

func (s *FindPurchaseSuite) TestError() {
	s.Run("should return not found when purchase does not exist", func() {
		t := s.T()

		repo := mocks.NewRepository(t)
		repo.EXPECT().FindByID(mock.Anything, testdata.ValidID).Return(purchasedomain.Purchase{}, purchasedomain.ErrPurchaseNotFound)

		h := findpurchasehandler.NewHandler(findpurchaseuc.New(repo))
		s.PrepareRouter(http.MethodGet, findpurchasehandler.Path, h.Handle)
		s.Serve(findpurchasehandler.Path, handler.WithPathParam("id", testdata.ValidID))

		resp := s.Response()
		require.Equal(t, http.StatusNotFound, resp.Code)
	})

	s.Run("should return bad request when id is invalid", func() {
		t := s.T()

		repo := mocks.NewRepository(t)
		h := findpurchasehandler.NewHandler(findpurchaseuc.New(repo))
		s.PrepareRouter(http.MethodGet, findpurchasehandler.Path, h.Handle)
		s.Serve(findpurchasehandler.Path, handler.WithPathParam("id", "abc"))

		resp := s.Response()
		require.Equal(t, http.StatusBadRequest, resp.Code)
	})

	s.Run("should return internal server error when use case fails", func() {
		t := s.T()

		repo := mocks.NewRepository(t)
		repo.EXPECT().FindByID(mock.Anything, testdata.ValidID).Return(purchasedomain.Purchase{}, errors.New("db error"))

		h := findpurchasehandler.NewHandler(findpurchaseuc.New(repo))
		s.PrepareRouter(http.MethodGet, findpurchasehandler.Path, h.Handle)
		s.Serve(findpurchasehandler.Path, handler.WithPathParam("id", testdata.ValidID))

		resp := s.Response()
		require.Equal(t, http.StatusInternalServerError, resp.Code)
	})
}
