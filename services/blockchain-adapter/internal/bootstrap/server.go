package bootstrap

import (
	grpcserver "github.com/cashback-platform/services/blockchain-adapter/internal/grpc"
	"go.uber.org/fx"
)

var Server = fx.Module("server",
	fx.Provide(grpcserver.NewTokenServer),
	fx.Invoke(grpcserver.StartServer),
)
