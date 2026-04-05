package bootstrap

import (
	"github.com/cashback-platform/services/mint-consumer/internal/infra/nats"

	"go.uber.org/fx"
)

var NATS = fx.Module("nats",
	fx.Provide(nats.NewClient),
)
