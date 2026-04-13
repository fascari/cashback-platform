//go:build integration

package repository_test

import (
	"context"
	"embed"
	"errors"
	"io/fs"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/cashback-platform/kit/testsuite"
	"github.com/cashback-platform/services/cashback-service-api/db/migrations"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/domain"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/repository"
	purchasetestdata "github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/repository/testdata"
)

//go:embed testdata/fixtures
var rawFixturesFS embed.FS

type RepositorySuite struct {
	testsuite.Suite
	repo repository.Repository
}

func TestRepositorySuite(t *testing.T) {
	suite.Run(t, new(RepositorySuite))
}

func (s *RepositorySuite) SetupSuite() {
	fixturesFS, err := fs.Sub(rawFixturesFS, "testdata/fixtures")
	s.Require().NoError(err)

	s.ConfigureDB()
	s.ConfigureFixtures(migrations.FS, fixturesFS)
	s.repo = repository.New(s.DB)
}

func (s *RepositorySuite) TestCreate() {
	s.MockNowFunc(purchasetestdata.FixedTime)
	ctx := context.Background()

	created, err := s.repo.Create(ctx, purchasetestdata.NewPurchase(purchasetestdata.UserID))

	s.Require().NoError(err)
	s.Require().Equal(purchasetestdata.CreatedPurchase(), created)
}

func (s *RepositorySuite) TestFindByID() {
	ctx := context.Background()

	found, err := s.repo.FindByID(ctx, purchasetestdata.PurchaseID)

	s.Require().NoError(err)
	s.Require().Equal(purchasetestdata.ExistingPurchase(), found)
}

func (s *RepositorySuite) TestFindByUserID() {
	ctx := context.Background()

	purchases, err := s.repo.FindByUserID(ctx, purchasetestdata.UserID)

	s.Require().NoError(err)
	s.Require().Len(purchases, 2)

	ids := []int64{purchases[0].ID, purchases[1].ID}
	s.Require().Contains(ids, purchasetestdata.PurchaseID)
	s.Require().Contains(ids, purchasetestdata.AnotherPurchaseID)
}

func (s *RepositorySuite) TestFindByUserID_Empty() {
	ctx := context.Background()

	purchases, err := s.repo.FindByUserID(ctx, purchasetestdata.NewUserID)

	s.Require().NoError(err)
	s.Require().Empty(purchases)
}

func (s *RepositorySuite) TestFindByID_NotFound() {
	ctx := context.Background()

	_, err := s.repo.FindByID(ctx, 999999)

	s.Require().True(errors.Is(err, domain.ErrPurchaseNotFound))
}
