package ethereum_test

import (
	"testing"

	"github.com/cashback-platform/kit/ethereum"
	"github.com/stretchr/testify/require"
)

func TestNewFromMnemonic_ShouldDeriveDeterministicAddress(t *testing.T) {
	const mnemonic = "test test test test test test test test test test test junk"
	const path = "m/44'/60'/0'/0/0"
	const wantAddress = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"

	wallet, err := ethereum.NewFromMnemonic(mnemonic, path)
	require.NoError(t, err)
	require.Equal(t, wantAddress, wallet.Address().Hex())
}

func TestNewFromMnemonic_ShouldFailOnInvalidMnemonic(t *testing.T) {
	_, err := ethereum.NewFromMnemonic("invalid mnemonic words here", "m/44'/60'/0'/0/0")
	require.Error(t, err)
}
