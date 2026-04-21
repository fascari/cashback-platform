//go:build e2e

package deposit_test

import (
	"context"
	"database/sql"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"

	e2esuite "github.com/cashback-platform/test/e2e/suite"
)

// mintABI contains only the mint function signature needed to trigger a Transfer event.
const mintABI = `[{"inputs":[{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"mint","outputs":[],"stateMutability":"nonpayable","type":"function"}]`

// Anvil deterministic deployer derived from the standard test mnemonic "test test…junk".
// The matching private key is well-known and safe to include in test code.
const deployerPrivateKey = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

// recipientAddress is the second Anvil test account; used only to receive minted tokens.
const recipientAddress = "0x70997970C51812dc3A010C7d01b50e0d17dc79C8"

type DepositSuite struct {
	suite.Suite
	db *sql.DB
}

func TestDepositSuite(t *testing.T) {
	if !e2esuite.BlockchainAvailable() {
		t.Skip("skipping deposit flow test: no EVM node reachable at ETHEREUM_RPC_URL")
	}
	suite.Run(t, new(DepositSuite))
}

func (s *DepositSuite) SetupSuite() {
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_DSN_BLOCKCHAIN"))
	s.Require().NoError(err)
	s.Require().NoError(db.Ping(), "blockchain DB must be reachable for deposit e2e tests")
	s.db = db
}

func (s *DepositSuite) TearDownSuite() {
	if s.db != nil {
		_ = s.db.Close()
	}
}

// TestDepositFlow_ShouldDetectOnChainTransfer mints tokens on Anvil and verifies that the
// deposit monitor detects the Transfer event and stores a row in detected_deposits.
func (s *DepositSuite) TestDepositFlow_ShouldDetectOnChainTransfer() {
	ctx := context.Background()

	rpcURL := os.Getenv("ETHEREUM_RPC_URL")
	if rpcURL == "" {
		rpcURL = "http://127.0.0.1:8545"
	}

	contractAddr := common.HexToAddress(os.Getenv("CONTRACT_ADDRESS"))

	client, err := ethclient.DialContext(ctx, rpcURL)
	s.Require().NoError(err)
	defer client.Close()

	privateKey, err := crypto.HexToECDSA(deployerPrivateKey)
	s.Require().NoError(err)

	chainID := big.NewInt(31337)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	s.Require().NoError(err)

	contractABI, err := abi.JSON(strings.NewReader(mintABI))
	s.Require().NoError(err)

	contract := bind.NewBoundContract(contractAddr, contractABI, client, client, client)

	recipient := common.HexToAddress(recipientAddress)
	amount := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil) // 1 token in wei

	tx, err := contract.Transact(auth, "mint", recipient, amount)
	s.Require().NoError(err, "mint transaction must be submitted to Anvil")

	_, err = bind.WaitMined(ctx, client, tx)
	s.Require().NoError(err, "mint transaction must be mined")

	txHash := tx.Hash().Hex()

	// The deposit monitor polls every 12 s by default; allow up to 45 s for detection.
	s.Require().Eventually(func() bool {
		var count int
		err := s.db.QueryRowContext(ctx,
			"SELECT COUNT(*) FROM detected_deposits WHERE transaction_hash = $1",
			txHash,
		).Scan(&count)
		return err == nil && count > 0
	}, 45*time.Second, 2*time.Second, "deposit %s not stored within 45s", txHash)
}
