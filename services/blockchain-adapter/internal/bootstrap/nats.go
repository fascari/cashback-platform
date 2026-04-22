package bootstrap

import (
	"go.uber.org/fx"

	infranats "github.com/cashback-platform/services/blockchain-adapter/internal/infra/nats"
)

var NATS = fx.Module("nats",
	fx.Provide(infranats.NewNATSClient),
)
