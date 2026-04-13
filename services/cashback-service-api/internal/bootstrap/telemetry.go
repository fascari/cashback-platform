package bootstrap

import (
	"context"

	"go.uber.org/fx"

	"github.com/cashback-platform/services/cashback-service-api/internal/config"
	"github.com/cashback-platform/services/cashback-service-api/internal/telemetry"
)

var Telemetry = fx.Module("telemetry",
	fx.Invoke(setupTelemetry),
)

type telemetryParams struct {
	fx.In

	LC  fx.Lifecycle
	Cfg config.Telemetry
}

func setupTelemetry(p telemetryParams) error {
	if !p.Cfg.Enabled {
		return nil
	}

	tp, err := telemetry.NewTracerProvider(context.Background(), serviceName, p.Cfg.OTLPEndpoint)
	if err != nil {
		return err
	}

	p.LC.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return tp.Shutdown(ctx)
		},
	})
	return nil
}
