package bootstrap

import (
	"github.com/cashback-platform/services/blockchain-adapter/internal/config"
	"go.uber.org/fx"
)

var Config = fx.Module("config",
	fx.Provide(config.NewConfig),
)
