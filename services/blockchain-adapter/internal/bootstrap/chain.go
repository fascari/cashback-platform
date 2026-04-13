package bootstrap

import (
	"go.uber.org/fx"

	"github.com/cashback-platform/services/blockchain-adapter/internal/chain"
	ethereuminfra "github.com/cashback-platform/services/blockchain-adapter/internal/infra/ethereum"
)

var Chain = fx.Module("chain",
	fx.Provide(
		ethereuminfra.New,
		newChainRegistry,
	),
)

func newChainRegistry(ethClient ethereuminfra.Client) *chain.Registry {
	return chain.NewRegistry(ethClient)
}
