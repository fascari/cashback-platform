package main

import (
	"go.uber.org/fx"

	"github.com/cashback-platform/services/blockchain-adapter/cmd/modules"
	"github.com/cashback-platform/services/blockchain-adapter/internal/bootstrap"
)

func main() {
	fx.New(
		bootstrap.Logger(),
		bootstrap.Config,
		bootstrap.Database,
		bootstrap.Redis,
		bootstrap.Ethereum,
		bootstrap.Server,
		modules.Token,
	).Run()
}
