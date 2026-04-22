//go:build integration

package repository_test

import (
	"context"
	"embed"
	"io/fs"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/cashback-platform/kit/testsuite"
	"github.com/cashback-platform/services/cashback-service-api/db/migrations"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/deposit/repository"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/deposit/repository/testdata"
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

func (s *RepositorySuite) TestSave_ShouldPersistDepositReceipt() {
	ctx := context.Background()

	receipt := testdata.NewDepositReceipt()
	receipt.TxHash = "0x1111111111111111111111111111111111111111111111111111111111111111ab"

	saved, err := s.repo.Save(ctx, receipt)

	s.Require().NoError(err)
	s.Require().Greater(saved.ID, int64(0))
	s.Require().Equal(receipt.UserID, saved.UserID)
	s.Require().Equal(receipt.TxHash, saved.TxHash)
	s.Require().Equal(receipt.FromAddress, saved.FromAddress)
	s.Require().Equal(receipt.Amount, saved.Amount)
	s.Require().Equal(receipt.ChainID, saved.ChainID)
	s.Require().Equal(receipt.BlockNumber, saved.BlockNumber)
	s.Require().Equal(receipt.DetectedAt, saved.DetectedAt)
	s.Require().False(saved.CreatedAt.IsZero())
}

func (s *RepositorySuite) TestExistsByTxHash_ShouldReturnTrueWhenExists() {
	ctx := context.Background()

	exists, err := s.repo.ExistsByTxHash(ctx, testdata.TxHash)

	s.Require().NoError(err)
	s.Require().True(exists)
}

func (s *RepositorySuite) TestExistsByTxHash_ShouldReturnFalseWhenNotExists() {
	ctx := context.Background()

	exists, err := s.repo.ExistsByTxHash(ctx, "0x0000000000000000000000000000000000000000000000000000000000000000ab")

	s.Require().NoError(err)
	s.Require().False(exists)
}
