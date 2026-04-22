package main

import (
	"github.com/cashback-platform/services/cashback-service-api/cmd/api/modules"
	"github.com/cashback-platform/services/cashback-service-api/internal/bootstrap"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		bootstrap.Logger(),
		bootstrap.Config,
		bootstrap.Database,
		bootstrap.NATS,
		bootstrap.GRPCClients,
		bootstrap.Outbox,
		bootstrap.Router,
		bootstrap.Server,
		bootstrap.Telemetry,
		modules.User,
		modules.Purchase,
		modules.Cashback,
		modules.Deposit,
	).Run()
}
