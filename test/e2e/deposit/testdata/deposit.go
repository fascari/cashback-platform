//go:build e2e

package testdata

import (
	_ "embed"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
)

//go:embed mint_abi.json
var MintABI string

const (
	// DeployerPrivateKey is the well-known Anvil account #0 private key.
	// Safe to embed in test code — it is a public test fixture, not a secret.
	DeployerPrivateKey = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

	recipientAddress = "0x70997970C51812dc3A010C7d01b50e0d17dc79C8"

	ChainID = int64(31337)
)

func ContractAddress() common.Address {
	return common.HexToAddress(os.Getenv("CONTRACT_ADDRESS"))
}

func RecipientAddr() common.Address {
	return common.HexToAddress(recipientAddress)
}

func MintAmount() *big.Int {
	return new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
}
