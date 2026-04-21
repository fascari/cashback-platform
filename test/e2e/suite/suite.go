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

const BaseURL = "http://localhost:18080"

type Suite struct {
	suite.Suite
	E               *httpexpect.Expect
	BlockchainDB    *sql.DB
	EthereumRPCURL  string
	db              *sql.DB
	loader          *testfixtures.Loader
}

func (s *Suite) SetupSuite() {
	resp, err := http.Get(BaseURL + "/health")
	if err != nil {
		s.T().Fatalf("API not reachable at %s — start services with mise run test:e2e", BaseURL)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		s.T().Fatalf("API not healthy at %s — status %d", BaseURL, resp.StatusCode)
	}

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_DSN_CASHBACK"))
	s.Require().NoError(err)
	s.Require().NoError(db.Ping())
	s.db = db

	if dsn := os.Getenv("POSTGRES_DSN_BLOCKCHAIN"); dsn != "" {
		bdb, err := sql.Open("postgres", dsn)
		s.Require().NoError(err)
		s.Require().NoError(bdb.Ping())
		s.BlockchainDB = bdb
	}

	rpcURL := os.Getenv("ETHEREUM_RPC_URL")
	if rpcURL == "" {
		rpcURL = "http://127.0.0.1:8545"
	}
	s.EthereumRPCURL = rpcURL
}

// ConfigureFixtures sets up the fixture loader for the given directory.
// The path is relative to the test package directory (where go test runs).
// Once configured, fixtures reload automatically before each test via SetupTest.
func (s *Suite) ConfigureFixtures(dir string) {
	loader, err := testfixtures.New(
		testfixtures.Database(s.db),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory(dir),
		testfixtures.DangerousSkipTestDatabaseCheck(),
	)
	s.Require().NoError(err)
	s.loader = loader
}

// SetupTest resets fixtures and reinitialises the HTTP client before each test.
func (s *Suite) SetupTest() {
	if s.loader != nil {
		s.Require().NoError(s.loader.Load())
	}

	s.E = httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  BaseURL + "/api/v1",
		Reporter: httpexpect.NewRequireReporter(s.T()),
		Client:   &http.Client{Timeout: 10 * time.Second},
	})
}

// TearDownSuite closes all database connections.
func (s *Suite) TearDownSuite() {
	if s.db != nil {
		_ = s.db.Close()
	}
	if s.BlockchainDB != nil {
		_ = s.BlockchainDB.Close()
	}
}

// BlockchainAvailable reports whether an EVM node is reachable.
// Falls back to http://127.0.0.1:8545 when ETHEREUM_RPC_URL is unset.
func BlockchainAvailable() bool {
	rpcURL := os.Getenv("ETHEREUM_RPC_URL")
	if rpcURL == "" {
		rpcURL = "http://127.0.0.1:8545"
	}
	c := &http.Client{Timeout: 2 * time.Second}
	resp, err := c.Post(rpcURL, "application/json", nil)
	if err != nil {
		return false
	}
	resp.Body.Close()
	return true
}
