//go:build integration

package repository_test

import (
	"context"
	"embed"
	"io/fs"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/cashback-platform/kit/testsuite"
	"github.com/cashback-platform/services/blockchain-adapter/db/migrations"
	"github.com/cashback-platform/services/blockchain-adapter/internal/app/deposit/repository"
	"github.com/cashback-platform/services/blockchain-adapter/internal/app/deposit/repository/testdata"
	depositassert "github.com/cashback-platform/services/blockchain-adapter/internal/app/deposit/repository/testdata/assert"
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

func (s *RepositorySuite) TestSave() {
	ctx := context.Background()
	deposit := testdata.NewDeposit()

	err := s.repo.Save(ctx, deposit)

	s.Require().NoError(err)
	depositassert.DepositSaved(s.T(), s.DB, deposit.ChainID, deposit.TransactionHash)
}

func (s *RepositorySuite) TestSave_Duplicate() {
	ctx := context.Background()

	err := s.repo.Save(ctx, testdata.ExistingDeposit())

	s.Require().NoError(err)
}

func (s *RepositorySuite) TestMaxBlockNumber() {
	ctx := context.Background()

	block, err := s.repo.MaxBlockNumber(ctx, testdata.ChainID)

	s.Require().NoError(err)
	s.Require().Equal(testdata.BlockNumber, block)
}

func (s *RepositorySuite) TestMaxBlockNumber_Empty() {
	ctx := context.Background()

	block, err := s.repo.MaxBlockNumber(ctx, "unknown-chain")

	s.Require().NoError(err)
	s.Require().Equal(int64(0), block)
}
