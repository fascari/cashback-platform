package bootstrap

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cashback-platform/kit/logger"
	"github.com/cashback-platform/services/cashback-service-api/internal/config"

	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
)

var Server = fx.Module("server",
	fx.Invoke(registerServer),
)

type serverParams struct {
	fx.In

	LC     fx.Lifecycle
	Router *chi.Mux `name:"main"`
	Cfg    config.Server
}

func registerServer(p serverParams) {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", p.Cfg.Port),
		Handler: p.Router,
	}

	p.LC.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() {
				logger.Info("Starting server", "port", p.Cfg.Port)
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Error("Server error", "error", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Shutting down server...")
			return server.Shutdown(ctx)
		},
	})
}
