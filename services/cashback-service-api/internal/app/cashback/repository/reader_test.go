//go:build integration

package repository_test

import (
	"context"
	"errors"

	cashdomain "github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/domain"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/repository/testdata"
)

func (s *RepositorySuite) TestFindByPurchaseID() {
	ctx := context.Background()

	found, err := s.repo.FindByPurchaseID(ctx, testdata.PurchaseID)

	s.Require().NoError(err)
	s.Require().Equal(testdata.PendingCashback(), found)
}

func (s *RepositorySuite) TestFindByUserID() {
	ctx := context.Background()

	cashbacks, err := s.repo.FindByUserID(ctx, testdata.UserID)

	s.Require().NoError(err)
	s.Require().ElementsMatch([]cashdomain.Cashback{
		testdata.PendingCashback(),
		testdata.MintedCashback(),
	}, cashbacks)
}

func (s *RepositorySuite) TestTotalByUserID() {
	ctx := context.Background()

	total, err := s.repo.TotalByUserID(ctx, testdata.UserID)

	s.Require().NoError(err)
	s.Require().Equal(10.0, total)
}

func (s *RepositorySuite) TestFindByID_NotFound() {
	ctx := context.Background()

	_, err := s.repo.FindByID(ctx, 999999)

	s.Require().True(errors.Is(err, cashdomain.ErrCashbackNotFound))
}

func (s *RepositorySuite) TestFindByPurchaseID_NotFound() {
	ctx := context.Background()

	_, err := s.repo.FindByPurchaseID(ctx, 999999)

	s.Require().True(errors.Is(err, cashdomain.ErrCashbackNotFound))
}
