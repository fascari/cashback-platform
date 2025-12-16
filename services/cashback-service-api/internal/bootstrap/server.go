package bootstrap

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cashback-platform/services/cashback-service-api/internal/config"
	"github.com/cashback-platform/services/cashback-service-api/pkg/logger"

	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
)

var Server = fx.Module("server",
	fx.Invoke(registerServer),
)

func registerServer(lc fx.Lifecycle, router *chi.Mux, cfg config.Server) {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: router,
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() {
				logger.Info("Starting server", "port", cfg.Port)
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
