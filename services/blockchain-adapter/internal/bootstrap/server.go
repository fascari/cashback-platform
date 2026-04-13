package bootstrap

import (
	"context"
	"fmt"
	"net"

	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/cashback-platform/kit/logger"
	"github.com/cashback-platform/services/blockchain-adapter/internal/config"
)

var Server = fx.Module("server",
	fx.Provide(newGRPCServer),
	fx.Invoke(startGRPCServer),
)

func newGRPCServer() *grpc.Server {
	server := grpc.NewServer()
	reflection.Register(server)
	return server
}

type serverParams struct {
	fx.In

	LC     fx.Lifecycle
	Server *grpc.Server
	Cfg    *config.Config
}

func startGRPCServer(p serverParams) {
	p.LC.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			listener, err := net.Listen("tcp", fmt.Sprintf(":%s", p.Cfg.GRPC.Port))
			if err != nil {
				return fmt.Errorf("listening on port %s: %w", p.Cfg.GRPC.Port, err)
			}

			go func() {
				logger.Info("gRPC server starting", "port", p.Cfg.GRPC.Port)
				if err := p.Server.Serve(listener); err != nil {
					logger.Error("gRPC server error", "error", err)
				}
			}()

			return nil
		},
		OnStop: func(_ context.Context) error {
			logger.Info("shutting down gRPC server")
			p.Server.GracefulStop()
			return nil
		},
	})
}
