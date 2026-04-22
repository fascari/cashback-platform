//go:build e2e

package suite

import (
	"database/sql"
	"net/http"
	"os"
	"time"

	"github.com/gavv/httpexpect/v2"
	"github.com/go-testfixtures/testfixtures/v3"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
)

const (
	BaseURL = "http://localhost:18080"

	envCashbackDSN    = "POSTGRES_DSN_CASHBACK"
	envBlockchainDSN  = "POSTGRES_DSN_BLOCKCHAIN"
	envEthereumRPCURL = "ETHEREUM_RPC_URL"

	defaultEthereumRPCURL = "http://127.0.0.1:8545"

	httpClientTimeout = 10 * time.Second
)

// Suite is the base for all e2e test suites. Embed it and override SetupSuite,
// calling s.Suite.SetupSuite() first.
type Suite struct {
	suite.Suite
	E              *httpexpect.Expect
	BlockchainDB   *sql.DB
	CashbackDB     *sql.DB
	EthereumRPCURL string
	loaders        []*testfixtures.Loader
}

func (s *Suite) SetupSuite() {
	s.requireAPIHealthy()
	s.CashbackDB = s.openDB(os.Getenv(envCashbackDSN))
	if dsn := os.Getenv(envBlockchainDSN); dsn != "" {
		s.BlockchainDB = s.openDB(dsn)
	}
	s.EthereumRPCURL = envOr(envEthereumRPCURL, defaultEthereumRPCURL)
}

// ConfigureFixtures sets up a fixture loader for the given database and directory.
// Call once per database in SetupSuite. Fixtures reload before each test via SetupTest.
// The directory path is relative to the test package directory (where go test runs).
func (s *Suite) ConfigureFixtures(db *sql.DB, dir string) {
	loader, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory(dir),
		testfixtures.DangerousSkipTestDatabaseCheck(),
	)
	s.Require().NoError(err)
	s.loaders = append(s.loaders, loader)
}

func (s *Suite) SetupTest() {
	for _, l := range s.loaders {
		s.Require().NoError(l.Load())
	}
	s.E = httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  BaseURL + "/api/v1",
		Reporter: httpexpect.NewRequireReporter(s.T()),
		Client:   &http.Client{Timeout: httpClientTimeout},
	})
}

func (s *Suite) TearDownSuite() {
	if s.CashbackDB != nil {
		_ = s.CashbackDB.Close()
	}
	if s.BlockchainDB != nil {
		_ = s.BlockchainDB.Close()
	}
}

// BlockchainAvailable reports whether an EVM node is reachable.
func BlockchainAvailable() bool {
	c := &http.Client{Timeout: 2 * time.Second}
	resp, err := c.Post(envOr(envEthereumRPCURL, defaultEthereumRPCURL), "application/json", nil)
	if err != nil {
		return false
	}
	_ = resp.Body.Close()
	return true
}

func (s *Suite) requireAPIHealthy() {
	resp, err := http.Get(BaseURL + "/health")
	if err != nil {
		s.T().Fatalf("API not reachable at %s — start services with mise run test:e2e", BaseURL)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		s.T().Fatalf("API not healthy at %s — status %d", BaseURL, resp.StatusCode)
	}
}

func (s *Suite) openDB(dsn string) *sql.DB {
	db, err := sql.Open("postgres", dsn)
	s.Require().NoError(err)
	s.Require().NoError(db.Ping())
	return db
}

// RowExists reports whether the given query returns at least one row.
func RowExists(db *sql.DB, query string, args ...any) bool {
	var count int
	return db.QueryRow(query, args...).Scan(&count) == nil && count > 0
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
