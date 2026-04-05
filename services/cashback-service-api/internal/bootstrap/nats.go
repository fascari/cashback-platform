package bootstrap

import (
	"github.com/cashback-platform/services/cashback-service-api/internal/infra/nats"

	"go.uber.org/fx"
)

// NATS provides the NATS client as an fx module.
var NATS = fx.Module("nats",
	fx.Provide(nats.NewNATSClient),
)
