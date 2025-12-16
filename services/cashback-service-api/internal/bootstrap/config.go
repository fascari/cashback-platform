package bootstrap

import (
	"github.com/cashback-platform/services/cashback-service-api/internal/config"

	"go.uber.org/fx"
)

var Config = fx.Module("config",
	fx.Provide(
		config.LoadDatabase,
		config.LoadNATS,
		config.LoadGRPC,
		config.LoadServer,
	),
)
