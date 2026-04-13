package calculatecashback_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/cashback-platform/kit/testsuite/handler"
	cashdomain "github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/domain"
	calculatecashbackhandler "github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/handler/calculatecashback"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/handler/calculatecashback/testdata"
	calculatecashbackuc "github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/usecase/calculatecashback"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/usecase/calculatecashback/mocks"
	purchasedomain "github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/domain"
	userdomain "github.com/cashback-platform/services/cashback-service-api/internal/app/user/domain"
)

type (
	CalculateCashbackSuite struct {
		handler.Suite
	}

	cashbackMocks struct {
		repo         *mocks.Repository
		purchaseRepo *mocks.PurchaseRepository
		userRepo     *mocks.UserRepository
	}
)

func TestCalculateCashback(t *testing.T) {
	suite.Run(t, &CalculateCashbackSuite{})
}

func newCashbackMocks(t *testing.T) cashbackMocks {
	return cashbackMocks{
		repo:         mocks.NewRepository(t),
		purchaseRepo: mocks.NewPurchaseRepository(t),
		userRepo:     mocks.NewUserRepository(t),
	}
}

func (m cashbackMocks) newHandler() calculatecashbackhandler.Handler {
	return calculatecashbackhandler.NewHandler(
		calculatecashbackuc.New(m.repo, m.purchaseRepo, m.userRepo),
	)
}

func (s *CalculateCashbackSuite) TestSuccess() {
	s.Run("should return created cashback", func() {
		t := s.T()
		m := newCashbackMocks(t)

		m.repo.EXPECT().FindByPurchaseID(mock.Anything, int64(1)).Return(cashdomain.Cashback{}, cashdomain.ErrCashbackNotFound)
		m.purchaseRepo.EXPECT().FindByID(mock.Anything, int64(1)).Return(purchasedomain.Purchase{ID: 1, UserID: 42, Amount: 100.0}, nil)
		m.userRepo.EXPECT().FindByID(mock.Anything, int64(42)).Return(userdomain.User{WalletAddress: "0xabc"}, nil)
		m.repo.EXPECT().CreateWithEvent(mock.Anything, mock.Anything, mock.Anything).Return(testdata.ExistingCashback(), nil)

		s.PrepareRouter(http.MethodPost, calculatecashbackhandler.Path, m.newHandler().Handle)
		s.Serve(calculatecashbackhandler.Path, handler.WithJSONBodyStruct(testdata.ValidPayload()))

		resp := s.Response()
		require.Equal(t, http.StatusCreated, resp.Code)
		require.Contains(t, resp.Body, `"status":"approved"`)
	})
}

func (s *CalculateCashbackSuite) TestError() {
	s.Run("should return bad request when payload is empty", func() {
		t := s.T()
		m := newCashbackMocks(t)

		s.PrepareRouter(http.MethodPost, calculatecashbackhandler.Path, m.newHandler().Handle)
		s.Serve(calculatecashbackhandler.Path, handler.WithJSONBodyStruct(testdata.MissingFieldsPayload()))

		resp := s.Response()
		require.Equal(t, http.StatusBadRequest, resp.Code)
		require.Contains(t, resp.Body, "purchase_id is a required field")
	})

	s.Run("should return bad request when body is invalid JSON", func() {
		t := s.T()
		m := newCashbackMocks(t)

		s.PrepareRouter(http.MethodPost, calculatecashbackhandler.Path, m.newHandler().Handle)
		s.Serve(calculatecashbackhandler.Path, handler.WithJSONBody("not valid json"))

		require.Equal(t, http.StatusBadRequest, s.Response().Code)
	})

	s.Run("should return internal server error when use case fails", func() {
		t := s.T()
		m := newCashbackMocks(t)

		m.repo.EXPECT().FindByPurchaseID(mock.Anything, mock.Anything).Return(cashdomain.Cashback{}, errors.New("db error"))

		s.PrepareRouter(http.MethodPost, calculatecashbackhandler.Path, m.newHandler().Handle)
		s.Serve(calculatecashbackhandler.Path, handler.WithJSONBodyStruct(testdata.ValidPayload()))

		require.Equal(t, http.StatusInternalServerError, s.Response().Code)
	})
}
