//go:build e2e

package cashback_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/suite"

	"github.com/cashback-platform/test/e2e/cashback/testdata"
	e2esuite "github.com/cashback-platform/test/e2e/suite"
)

type CashbackSuite struct {
	e2esuite.Suite
}

func TestCashbackSuite(t *testing.T) {
	suite.Run(t, new(CashbackSuite))
}

func (s *CashbackSuite) SetupSuite() {
	s.Suite.SetupSuite()
	s.Suite.ConfigureFixtures("testdata/fixtures")
}

func (s *CashbackSuite) TestCashbackFlow_ShouldIncrementBalanceAfterMint() {
	if !e2esuite.BlockchainAvailable() {
		s.T().Skip("skipping blockchain flow test: no EVM node reachable at ETHEREUM_RPC_URL")
	}

	s.E.POST("/cashback/calculate").WithJSON(map[string]any{
		"purchase_id": testdata.FlowPurchaseID,
	}).Expect().Status(http.StatusCreated).
		JSON().Object().HasValue("status", "approved")

	balanceBefore := s.E.GET(fmt.Sprintf("/users/%s/balance", testdata.FlowUserID)).
		Expect().Status(http.StatusOK).
		JSON().Object().Value("balance").String().Raw()

	// polling uses a non-fatal reporter so Eventually can retry on transient errors.
	pollExpect := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  e2esuite.BaseURL + "/api/v1",
		Reporter: httpexpect.NewAssertReporter(s.T()),
		Client:   &http.Client{Timeout: 15 * time.Second},
	})

	var balanceAfter string
	s.Require().Eventually(func() bool {
		balanceAfter = pollExpect.GET(fmt.Sprintf("/users/%s/balance", testdata.FlowUserID)).
			Expect().Status(http.StatusOK).
			JSON().Object().Value("balance").String().Raw()
		return balanceAfter != balanceBefore
	}, 60*time.Second, 3*time.Second)
}

func (s *CashbackSuite) TestCalculateCashback_ShouldReturnConflictWhenAlreadyCalculated() {
	s.E.POST("/cashback/calculate").WithJSON(map[string]any{
		"purchase_id": testdata.IdemPurchaseID,
	}).Expect().Status(http.StatusCreated)

	s.E.POST("/cashback/calculate").WithJSON(map[string]any{
		"purchase_id": testdata.IdemPurchaseID,
	}).Expect().Status(http.StatusConflict)
}

func (s *CashbackSuite) TestCalculateCashback_ShouldReturnNotFoundForUnknownPurchase() {
	s.E.POST("/cashback/calculate").WithJSON(map[string]any{
		"purchase_id": "999999",
	}).Expect().Status(http.StatusNotFound)
}
