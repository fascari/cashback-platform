package bootstrap

import (
	"github.com/cashback-platform/services/cashback-service-api/internal/infra/grpc"

	"go.uber.org/fx"
)

var GRPCClients = fx.Module("grpc-clients",
	fx.Provide(grpc.NewBlockchainAdapterClient),
)
