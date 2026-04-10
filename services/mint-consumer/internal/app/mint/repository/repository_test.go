//go:build integration

package repository_test

import (
	"context"
	"embed"
	"io/fs"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/cashback-platform/kit/testsuite"
	"github.com/cashback-platform/services/mint-consumer/db/migrations"
	"github.com/cashback-platform/services/mint-consumer/internal/app/mint/domain"
	"github.com/cashback-platform/services/mint-consumer/internal/app/mint/repository"
	"github.com/cashback-platform/services/mint-consumer/internal/app/mint/repository/testdata"
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

func (s *RepositorySuite) TestCreateMintRequest() {
	ctx := context.Background()

	req, err := s.repo.CreateMintRequest(ctx, testdata.NewPendingMintRequest())

	s.Require().NoError(err)
	s.Require().Greater(req.ID, int64(0))
	s.Require().Equal(domain.MintRequestStatusPending, req.Status)
}

func (s *RepositorySuite) TestExistsProcessedEvent_NotFound() {
	ctx := context.Background()

	exists, err := s.repo.ExistsProcessedEvent(ctx, testdata.EventID)

	s.Require().NoError(err)
	s.Require().False(exists)
}

func (s *RepositorySuite) TestCreateProcessedEvent() {
	ctx := context.Background()

	err := s.repo.CreateProcessedEvent(ctx, testdata.EventID, "cashback.approved")
	s.Require().NoError(err)

	exists, err := s.repo.ExistsProcessedEvent(ctx, testdata.EventID)
	s.Require().NoError(err)
	s.Require().True(exists)
}

func (s *RepositorySuite) TestCreateProcessedEvent_Duplicate() {
	ctx := context.Background()

	err := s.repo.CreateProcessedEvent(ctx, testdata.EventID, "cashback.approved")
	s.Require().NoError(err)

	err = s.repo.CreateProcessedEvent(ctx, testdata.EventID, "cashback.approved")
	s.Require().Error(err)
}

func (s *RepositorySuite) TestMarkCompleted() {
	ctx := context.Background()

	req, err := s.repo.CreateMintRequest(ctx, testdata.NewPendingMintRequest())
	s.Require().NoError(err)

	err = s.repo.MarkCompleted(ctx, req.ID, "0xabc123", 9999)
	s.Require().NoError(err)

	var status string
	s.DB.Raw("SELECT status FROM mint.mint_requests WHERE id = ?", req.ID).Scan(&status)
	s.Require().Equal("completed", status)

	var txHash string
	s.DB.Raw("SELECT transaction_hash FROM mint.mint_requests WHERE id = ?", req.ID).Scan(&txHash)
	s.Require().Equal("0xabc123", txHash)

	var blockNumber int64
	s.DB.Raw("SELECT block_number FROM mint.mint_requests WHERE id = ?", req.ID).Scan(&blockNumber)
	s.Require().Equal(int64(9999), blockNumber)
}

func (s *RepositorySuite) TestMarkFailed() {
	ctx := context.Background()

	req, err := s.repo.CreateMintRequest(ctx, testdata.NewPendingMintRequest())
	s.Require().NoError(err)

	nextRetry := time.Now().UTC().Add(1 * time.Minute)
	err = s.repo.MarkFailed(ctx, req.ID, "timeout", "rpc timeout", &nextRetry)
	s.Require().NoError(err)

	var status string
	s.DB.Raw("SELECT status FROM mint.mint_requests WHERE id = ?", req.ID).Scan(&status)
	s.Require().Equal("failed", status)

	var errorCode string
	s.DB.Raw("SELECT error_code FROM mint.mint_requests WHERE id = ?", req.ID).Scan(&errorCode)
	s.Require().Equal("timeout", errorCode)

	var retryCount int
	s.DB.Raw("SELECT retry_count FROM mint.mint_requests WHERE id = ?", req.ID).Scan(&retryCount)
	s.Require().Equal(1, retryCount)
}

func (s *RepositorySuite) TestFindFailedRetryable() {
	ctx := context.Background()

	results, err := s.repo.FindFailedRetryable(ctx, 10)

	s.Require().NoError(err)
	s.Require().Len(results, 1)
	s.Require().Equal(int64(1), results[0].ID)
}
