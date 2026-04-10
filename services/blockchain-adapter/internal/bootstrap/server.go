package bootstrap

import (
	"github.com/cashback-platform/services/blockchain-adapter/internal/app/token/handler"
	"go.uber.org/fx"
)

var Server = fx.Module("server",
	fx.Provide(handler.NewTokenServer),
	fx.Invoke(handler.StartServer),
)
