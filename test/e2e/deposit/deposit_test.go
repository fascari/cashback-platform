//go:build e2e

package deposit_test

import (
	"context"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
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

func (s *DepositSuite) TestDepositFlow_ShouldDetectOnChainTransfer() {
	ctx := context.Background()

	client, err := ethclient.DialContext(ctx, s.EthereumRPCURL)
	s.Require().NoError(err)
	defer client.Close()

	privateKey, err := crypto.HexToECDSA(testdata.DeployerPrivateKey)
	s.Require().NoError(err)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(testdata.ChainID))
	s.Require().NoError(err)

	contractABI, err := abi.JSON(strings.NewReader(testdata.MintABI))
	s.Require().NoError(err)

	contract := bind.NewBoundContract(testdata.ContractAddress(), contractABI, client, client, client)

	tx, err := contract.Transact(auth, "mint", testdata.RecipientAddr(), testdata.MintAmount())
	s.Require().NoError(err, "mint transaction must be submitted to Anvil")

	_, err = bind.WaitMined(ctx, client, tx)
	s.Require().NoError(err, "mint transaction must be mined")

	txHash := tx.Hash().Hex()

	s.Require().Eventually(func() bool {
		return depositassert.IsDetected(s.BlockchainDB, txHash)
	}, 45*time.Second, 2*time.Second, "deposit %s not stored within 45s", txHash)
}
