package bootstrap

import (
	infragrpc "github.com/cashback-platform/services/mint-consumer/internal/infra/grpc"

	"go.uber.org/fx"
)

var GRPCClients = fx.Module("grpc-clients",
	fx.Provide(infragrpc.NewBlockchainAdapterClient),
)
