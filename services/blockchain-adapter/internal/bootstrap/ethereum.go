package bootstrap

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/fx"

	"github.com/cashback-platform/services/blockchain-adapter/internal/config"
	"github.com/cashback-platform/services/blockchain-adapter/internal/contracts"
	ethereumpkg "github.com/cashback-platform/services/blockchain-adapter/pkg/ethereum"
)

var Ethereum = fx.Module("ethereum",
	fx.Provide(
		newEthereumClient,
		newWallet,
		newCashbackToken,
	),
)

func newEthereumClient(cfg *config.Config) (*ethereumpkg.Client, error) {
	return ethereumpkg.New(cfg.Ethereum.RPCURL)
}

func newWallet(cfg *config.Config) (*ethereumpkg.Wallet, error) {
	return ethereumpkg.NewFromMnemonic(cfg.Wallet.Mnemonic, cfg.Wallet.DerivationPath)
}

func newCashbackToken(cfg *config.Config, client *ethereumpkg.Client) (*contracts.CashbackToken, error) {
	addr := common.HexToAddress(cfg.Ethereum.ContractAddress)
	token, err := contracts.NewCashbackToken(addr, client.Inner())
	if err != nil {
		return nil, fmt.Errorf("instantiate cashback token contract: %w", err)
	}
	return token, nil
}
