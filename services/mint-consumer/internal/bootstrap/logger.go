package bootstrap

import (
	"os"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	"github.com/cashback-platform/kit/logger"
)

func init() {
	logger.Init()
}

// Logger returns an Fx option that configures the FX event logger to write to stderr.
func Logger() fx.Option {
	return fx.WithLogger(func() fxevent.Logger {
		return &fxevent.ConsoleLogger{W: os.Stderr}
	})
}
