//go:build integration

package repository_test

import (
	"context"

	cashdomain "github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/domain"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/repository/testdata"
)

func (s *RepositorySuite) TestCreate() {
	s.MockNowFunc(testdata.FixedTime)
	ctx := context.Background()

	created, err := s.repo.Create(ctx, testdata.NewCashback())

	s.Require().NoError(err)
	s.Require().Equal(testdata.CreatedCashback(), created)
}

func (s *RepositorySuite) TestCreateWithEvent() {
	s.MockNowFunc(testdata.FixedTime)
	ctx := context.Background()

	created, err := s.repo.CreateWithEvent(ctx, testdata.NewCashback(), func(c cashdomain.Cashback) any {
		return map[string]any{"cashback_id": c.ID, "amount": c.Amount}
	})

	s.Require().NoError(err)
	s.Require().Equal(testdata.CreatedCashback(), created)
}
