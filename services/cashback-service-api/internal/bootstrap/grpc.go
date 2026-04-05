package bootstrap

import (
	"github.com/cashback-platform/services/cashback-service-api/internal/infra/grpc"

	"go.uber.org/fx"
)

// GRPCClients provides gRPC clients as an fx module.
var GRPCClients = fx.Module("grpc-clients",
	fx.Provide(grpc.NewBlockchainAdapterClient),
)
