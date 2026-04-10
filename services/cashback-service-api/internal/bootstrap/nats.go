package bootstrap

import (
	"github.com/cashback-platform/services/cashback-service-api/internal/infra/nats"

	"go.uber.org/fx"
)

var NATS = fx.Module("nats",
	fx.Provide(nats.NewNATSClient),
)
