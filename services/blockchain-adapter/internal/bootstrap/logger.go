package bootstrap

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	"github.com/cashback-platform/kit/logger"
)

func init() {
	logger.Init()
}

func Logger() fx.Option {
	return fx.WithLogger(func() fxevent.Logger {
		return &fxevent.ZapLogger{
			Logger: logger.Zap().WithOptions(zap.IncreaseLevel(zap.ErrorLevel)),
		}
	})
}
