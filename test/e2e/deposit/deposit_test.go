//go:build e2e

package deposit_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/cashback-platform/test/e2e/deposit/testdata"
	depositassert "github.com/cashback-platform/test/e2e/deposit/testdata/assert"
	e2esuite "github.com/cashback-platform/test/e2e/suite"
)

type DepositSuite struct {
	e2esuite.Suite
}

func TestDepositSuite(t *testing.T) {
	if !e2esuite.BlockchainAvailable() {
		t.Skip("skipping deposit flow test: no EVM node reachable at ETHEREUM_RPC_URL")
	}
	suite.Run(t, new(DepositSuite))
}

func (s *DepositSuite) SetupSuite() {
	s.Suite.SetupSuite()
	s.Suite.ConfigureFixtures(s.CashbackDB, "testdata/fixtures")
}

func (s *DepositSuite) TestDepositFlow_ShouldDetectOnChainTransfer() {
	txHash := testdata.MintTokens(s.T(), s.EthereumRPCURL)

	s.Require().Eventually(func() bool {
		return depositassert.IsDetected(s.BlockchainDB, txHash)
	}, 45*time.Second, 2*time.Second, "deposit %s not stored within 45s", txHash)
}

func (s *DepositSuite) TestDepositFlow_ShouldCreateDepositReceiptAndCashback() {
	testdata.MintTokens(s.T(), s.EthereumRPCURL)
	txHash := testdata.DepositTokens(s.T(), s.EthereumRPCURL)

	s.Require().Eventually(func() bool {
		return depositassert.IsDetected(s.BlockchainDB, txHash)
	}, 45*time.Second, 2*time.Second, "deposit %s not stored within 45s", txHash)

	s.Require().Eventually(func() bool {
		return depositassert.DepositReceiptExists(s.CashbackDB, txHash)
	}, 30*time.Second, 2*time.Second, "deposit_receipt for tx_hash %s not found", txHash)

	s.Require().Eventually(func() bool {
		return depositassert.CashbackFromDepositExists(s.CashbackDB, testdata.RecipientWalletAddress)
	}, 30*time.Second, 2*time.Second, "cashback from deposit for wallet %s not found", testdata.RecipientWalletAddress)
}
