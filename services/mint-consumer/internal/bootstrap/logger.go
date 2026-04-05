package bootstrap

import (
	"github.com/cashback-platform/kit/logger"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func init() {
	logger.Init()
}

// Logger returns an fx.Option that silences the default FX console logger.
func Logger() fx.Option {
	return fx.WithLogger(func() fxevent.Logger {
		return &fxevent.ConsoleLogger{W: nil}
	})
}
