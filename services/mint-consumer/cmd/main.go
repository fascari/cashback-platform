package main

import (
	"github.com/cashback-platform/services/mint-consumer/cmd/modules"
	"github.com/cashback-platform/services/mint-consumer/internal/bootstrap"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		bootstrap.Config,
		bootstrap.Database,
		bootstrap.NATS,
		bootstrap.GRPCClients,
		modules.Mint,
	).Run()
}
