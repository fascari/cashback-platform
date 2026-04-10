package createpurchase_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/cashback-platform/kit/testsuite/handler"
	purchasedomain "github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/domain"
	createpurchasehandler "github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/handler/createpurchase"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/handler/createpurchase/testdata"
	createpurchaseuc "github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/usecase/createpurchase"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/usecase/createpurchase/mocks"
)

type CreatePurchaseSuite struct {
	handler.Suite
}

func TestCreatePurchase(t *testing.T) {
	suite.Run(t, &CreatePurchaseSuite{})
}

func (s *CreatePurchaseSuite) TestSuccess() {
	s.Run("should return 201 with purchase", func() {
		t := s.T()

		repo := mocks.NewRepository(t)
		repo.EXPECT().Create(mock.Anything, mock.Anything).Return(testdata.CreatedPurchase(), nil)

		h := createpurchasehandler.NewHandler(createpurchaseuc.New(repo))
		s.PrepareRouter(http.MethodPost, createpurchasehandler.Path, h.Handle)
		s.Serve(createpurchasehandler.Path, handler.WithJSONBodyStruct(testdata.ValidPayload()))

		resp := s.Response()
		require.Equal(t, http.StatusCreated, resp.Code)
		require.Contains(t, resp.Body, `"status":"pending"`)
	})
}

func (s *CreatePurchaseSuite) TestError() {
	s.Run("should return bad request when payload is empty", func() {
		t := s.T()

		repo := mocks.NewRepository(t)
		h := createpurchasehandler.NewHandler(createpurchaseuc.New(repo))
		s.PrepareRouter(http.MethodPost, createpurchasehandler.Path, h.Handle)
		s.Serve(createpurchasehandler.Path, handler.WithJSONBodyStruct(testdata.MissingFieldsPayload()))

		resp := s.Response()
		require.Equal(t, http.StatusBadRequest, resp.Code)
		require.Contains(t, resp.Body, "user_id is a required field")
		require.Contains(t, resp.Body, "merchant is a required field")
	})

	s.Run("should return bad request when amount is zero", func() {
		t := s.T()

		repo := mocks.NewRepository(t)
		h := createpurchasehandler.NewHandler(createpurchaseuc.New(repo))
		s.PrepareRouter(http.MethodPost, createpurchasehandler.Path, h.Handle)
		s.Serve(createpurchasehandler.Path, handler.WithJSONBodyStruct(testdata.ZeroAmountPayload()))

		resp := s.Response()
		require.Equal(t, http.StatusBadRequest, resp.Code)
		require.Contains(t, resp.Body, "amount must be greater than 0")
	})

	s.Run("should return internal server error when use case fails", func() {
		t := s.T()

		repo := mocks.NewRepository(t)
		repo.EXPECT().Create(mock.Anything, mock.Anything).Return(purchasedomain.Purchase{}, errors.New("db error"))

		h := createpurchasehandler.NewHandler(createpurchaseuc.New(repo))
		s.PrepareRouter(http.MethodPost, createpurchasehandler.Path, h.Handle)
		s.Serve(createpurchasehandler.Path, handler.WithJSONBodyStruct(testdata.ValidPayload()))

		resp := s.Response()
		require.Equal(t, http.StatusInternalServerError, resp.Code)
	})

	s.Run("should return bad request when body is invalid JSON", func() {
		t := s.T()

		repo := mocks.NewRepository(t)
		h := createpurchasehandler.NewHandler(createpurchaseuc.New(repo))
		s.PrepareRouter(http.MethodPost, createpurchasehandler.Path, h.Handle)
		s.Serve(createpurchasehandler.Path, handler.WithJSONBody("not valid json"))

		resp := s.Response()
		require.Equal(t, http.StatusBadRequest, resp.Code)
	})
}
