//go:build e2e

package testdata

import (
	"context"
	_ "embed"
	"math/big"
	"os"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

//go:embed mint_abi.json
var MintABI string

const (
	// DeployerPrivateKey is the well-known Anvil account #0 private key.
	// Safe to embed in test code — it is a public test fixture, not a secret.
	DeployerPrivateKey = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

	// RecipientPrivateKey is the well-known Anvil account #1 private key.
	// Safe to embed in test code — it is a public test fixture, not a secret.
	RecipientPrivateKey = "59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d"

	RecipientWalletAddress = "0x70997970C51812dc3A010C7d01b50e0d17dc79C8"

	ChainID = int64(31337)
)

func ContractAddress() common.Address {
	return common.HexToAddress(os.Getenv("CONTRACT_ADDRESS"))
}

func DeployerAddr() common.Address {
	return common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")
}

func RecipientAddr() common.Address {
	return common.HexToAddress(RecipientWalletAddress)
}

func MintAmount() *big.Int {
	return new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
}

// MintTokens submits and mines a mint transaction on the local Anvil node,
// returning the transaction hash.
func MintTokens(t *testing.T, rpcURL string) string {
	t.Helper()
	ctx := context.Background()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		t.Fatalf("dial ethereum client: %v", err)
	}
	defer client.Close()

	privateKey, err := crypto.HexToECDSA(DeployerPrivateKey)
	if err != nil {
		t.Fatalf("parse private key: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(ChainID))
	if err != nil {
		t.Fatalf("new transactor: %v", err)
	}

	contractABI, err := abi.JSON(strings.NewReader(MintABI))
	if err != nil {
		t.Fatalf("parse ABI: %v", err)
	}

	contract := bind.NewBoundContract(ContractAddress(), contractABI, client, client, client)

	tx, err := contract.Transact(auth, "mint", RecipientAddr(), MintAmount())
	if err != nil {
		t.Fatalf("mint transaction: %v", err)
	}

	_, err = bind.WaitMined(ctx, client, tx)
	if err != nil {
		t.Fatalf("wait mined: %v", err)
	}

	return tx.Hash().Hex()
}

// DepositTokens transfers tokens from RecipientAddr (Anvil #1) to DeployerAddr,
// simulating a user depositing tokens to the platform. MintTokens must be called
// first to fund RecipientAddr. Returns the transaction hash of the transfer.
func DepositTokens(t *testing.T, rpcURL string) string {
	t.Helper()
	ctx := context.Background()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		t.Fatalf("dial ethereum client: %v", err)
	}
	defer client.Close()

	privateKey, err := crypto.HexToECDSA(RecipientPrivateKey)
	if err != nil {
		t.Fatalf("parse private key: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(ChainID))
	if err != nil {
		t.Fatalf("new transactor: %v", err)
	}

	contractABI, err := abi.JSON(strings.NewReader(MintABI))
	if err != nil {
		t.Fatalf("parse ABI: %v", err)
	}

	contract := bind.NewBoundContract(ContractAddress(), contractABI, client, client, client)

	tx, err := contract.Transact(auth, "transfer", DeployerAddr(), MintAmount())
	if err != nil {
		t.Fatalf("transfer transaction: %v", err)
	}

	_, err = bind.WaitMined(ctx, client, tx)
	if err != nil {
		t.Fatalf("wait mined: %v", err)
	}

	return tx.Hash().Hex()
}
