//go:build integration

package testsuite

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"time"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	dbImage    = "postgres:15-alpine"
	dbUser     = "test"
	dbPassword = "test"
	dbName     = "testdb"
	dbPort     = "5432/tcp"
	dbTimeout  = 30 * time.Second
)

// Suite provides a reusable database container, GORM connection,
// and fixture loader for repository integration tests.
type Suite struct {
	suite.Suite
	DB        *gorm.DB
	container testcontainers.Container
	loader    *testfixtures.Loader
}

// ConfigureDB spins up a postgres:15-alpine container and opens a GORM connection.
// Returns the DSN for further use if needed.
func (s *Suite) ConfigureDB() string {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        dbImage,
		ExposedPorts: []string{dbPort},
		Env: map[string]string{
			"POSTGRES_USER":     dbUser,
			"POSTGRES_PASSWORD": dbPassword,
			"POSTGRES_DB":       dbName,
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).
			WithStartupTimeout(dbTimeout),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	s.Require().NoError(err)
	s.container = container

	host, err := container.Host(ctx)
	s.Require().NoError(err)
	port, err := container.MappedPort(ctx, "5432")
	s.Require().NoError(err)

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port.Port(), dbUser, dbPassword, dbName)
	db, err := gorm.Open(gormpostgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	s.Require().NoError(err)
	s.DB = db

	return dsn
}

// ConfigureFixtures runs migrations and initialises the fixture loader.
//
// migrationsFS is an embedded FS containing the SQL migration files (*.up.sql, *.down.sql).
// fixturesFS is an embedded FS containing the YAML fixture files loaded before each test.
func (s *Suite) ConfigureFixtures(migrationsFS fs.FS, fixturesFS fs.FS) {
	sqlDB, err := s.DB.DB()
	s.Require().NoError(err)

	src, err := iofs.New(migrationsFS, ".")
	s.Require().NoError(err)

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	s.Require().NoError(err)

	m, err := migrate.NewWithInstance("iofs", src, "postgres", driver)
	s.Require().NoError(err)

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		s.Require().NoError(err)
	}

	s.loader, err = testfixtures.New(
		testfixtures.Database(sqlDB),
		testfixtures.Dialect("postgres"),
		testfixtures.FS(fixturesFS),
		testfixtures.Directory("."),
	)
	s.Require().NoError(err)
}

// TearDownSuite terminates the PostgreSQL container.
func (s *Suite) TearDownSuite() {
	if s.container != nil {
		_ = s.container.Terminate(context.Background())
	}
}

// SetupTest reloads fixtures before each test.
func (s *Suite) SetupTest() {
	s.Require().NoError(s.loader.Load())
}
