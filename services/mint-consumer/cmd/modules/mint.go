package modules

import (
	"github.com/cashback-platform/services/mint-consumer/internal/app/mint/handler"
	"github.com/cashback-platform/services/mint-consumer/internal/app/mint/handler/cashbackapproved"
	"github.com/cashback-platform/services/mint-consumer/internal/app/mint/repository"
	"github.com/cashback-platform/services/mint-consumer/internal/app/mint/usecase/mintcashback"
	"github.com/cashback-platform/services/mint-consumer/internal/app/mint/usecase/retrymints"
	infragrpc "github.com/cashback-platform/services/mint-consumer/internal/infra/grpc"
	"go.uber.org/fx"
)

var (
	mintFactories = fx.Provide(
		repository.New,
		mintcashback.NewUseCase,
		retrymints.NewUseCase,
		cashbackapproved.New,
		handler.NewCashback,
	)

	mintDependencies = fx.Provide(
		func(r repository.Repository) mintcashback.Repository { return r },
		func(r repository.Repository) retrymints.Repository { return r },
		func(c *infragrpc.Client) mintcashback.BlockchainClient { return c },
		func(c *infragrpc.Client) retrymints.BlockchainClient { return c },
	)

	mintInvokes = fx.Invoke(handler.StartConsumer)

	Mint = fx.Options(
		mintFactories,
		mintDependencies,
		mintInvokes,
	)
)
