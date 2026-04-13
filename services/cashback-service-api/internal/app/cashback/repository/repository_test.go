//go:build integration

package repository_test

import (
	"embed"
	"io/fs"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/cashback-platform/kit/testsuite"
	"github.com/cashback-platform/services/cashback-service-api/db/migrations"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/repository"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/repository/mocks"
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

	outboxWriter := mocks.NewOutboxWriter(s.T())
	outboxWriter.EXPECT().CreateWithTx(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	s.repo = repository.New(s.DB, outboxWriter)
}
