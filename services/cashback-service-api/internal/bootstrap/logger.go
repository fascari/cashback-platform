package bootstrap

import (
	"github.com/cashback-platform/services/cashback-service-api/pkg/logger"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func init() {
	logger.Init()
}

func Logger() fx.Option {
	return fx.WithLogger(func() fxevent.Logger {
		return &fxevent.ConsoleLogger{W: nil}
	})
}
