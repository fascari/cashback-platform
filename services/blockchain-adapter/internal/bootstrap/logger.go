package bootstrap

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	"github.com/cashback-platform/kit/logger"
)

func init() {
	logger.Init()
}

// Logger returns an Fx option that suppresses Fx internal event logs.
func Logger() fx.Option {
	return fx.WithLogger(func() fxevent.Logger {
		return fxevent.NopLogger
	})
}
