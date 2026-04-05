package bootstrap

import (
	"github.com/cashback-platform/services/mint-consumer/internal/config"

	"go.uber.org/fx"
)

var Config = fx.Module("config",
	fx.Provide(config.LoadDatabase),
	fx.Provide(config.LoadNATS),
	fx.Provide(config.LoadGRPC),
)
