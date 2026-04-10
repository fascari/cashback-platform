package testdata

import (
	"testing"

	"github.com/stretchr/testify/require"

	ethereumpkg "github.com/cashback-platform/kit/ethereum"
	"github.com/cashback-platform/services/blockchain-adapter/internal/config"
)

const (
	mnemonic       = "test test test test test test test test test test test junk"
	derivationPath = "m/44'/60'/0'/0/0"
)

func Wallet(t *testing.T) *ethereumpkg.Wallet {
	t.Helper()
	w, err := ethereumpkg.NewFromMnemonic(mnemonic, derivationPath)
	require.NoError(t, err)
	return w
}

func Config() *config.Config {
	return &config.Config{
		Ethereum: config.EthereumConfig{
			ChainID:         11155111,
			ContractAddress: "0x0000000000000000000000000000000000000000",
		},
	}
}
