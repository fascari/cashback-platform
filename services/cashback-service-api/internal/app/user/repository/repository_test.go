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
	"github.com/cashback-platform/services/cashback-service-api/internal/app/user/domain"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/user/repository"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/user/repository/testdata"
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
	s.MockNowFunc(testdata.FixedTime)
	ctx := context.Background()

	created, err := s.repo.Create(ctx, testdata.AnotherUser())

	s.Require().NoError(err)
	s.Require().Equal(testdata.CreatedUser(), created)
}

func (s *RepositorySuite) TestFindByID() {
	ctx := context.Background()

	found, err := s.repo.FindByID(ctx, testdata.UserID)

	s.Require().NoError(err)
	s.Require().Equal(testdata.ExistingUser(), found)
}

func (s *RepositorySuite) TestFindByExternalID() {
	ctx := context.Background()

	found, err := s.repo.FindByExternalID(ctx, testdata.ExistingUser().ExternalID)

	s.Require().NoError(err)
	s.Require().Equal(testdata.ExistingUser(), found)
}

func (s *RepositorySuite) TestFindByEmail() {
	ctx := context.Background()

	found, err := s.repo.FindByEmail(ctx, testdata.ExistingUser().Email)

	s.Require().NoError(err)
	s.Require().Equal(testdata.ExistingUser(), found)
}

func (s *RepositorySuite) TestFindByID_NotFound() {
	ctx := context.Background()

	_, err := s.repo.FindByID(ctx, 999999)

	s.Require().True(errors.Is(err, domain.ErrUserNotFound))
}

func (s *RepositorySuite) TestFindByExternalID_NotFound() {
	ctx := context.Background()

	_, err := s.repo.FindByExternalID(ctx, "ghost-ext")

	s.Require().True(errors.Is(err, domain.ErrUserNotFound))
}

func (s *RepositorySuite) TestFindByEmail_NotFound() {
	ctx := context.Background()

	_, err := s.repo.FindByEmail(ctx, "ghost@example.com")

	s.Require().True(errors.Is(err, domain.ErrUserNotFound))
}
