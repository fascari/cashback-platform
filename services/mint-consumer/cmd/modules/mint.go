package modules

import (
	"github.com/cashback-platform/kit/gormtx"
	"github.com/cashback-platform/services/mint-consumer/internal/app/mint/repository"
	"github.com/cashback-platform/services/mint-consumer/internal/app/mint/usecase/processcashbackapproved"
	"github.com/cashback-platform/services/mint-consumer/internal/app/mint/usecase/retryfailedmints"
	"github.com/cashback-platform/services/mint-consumer/internal/consumer"
	infragrpc "github.com/cashback-platform/services/mint-consumer/internal/infra/grpc"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var (
	mintFactories = fx.Provide(
		repository.New,
		processcashbackapproved.NewUseCase,
		retryfailedmints.NewUseCase,
		consumer.NewCashback,
	)

	mintDependencies = fx.Provide(
		func(r repository.Repository) processcashbackapproved.Repository { return r },
		func(r repository.Repository) retryfailedmints.Repository { return r },
		func(c *infragrpc.Client) processcashbackapproved.BlockchainClient { return c },
		func(c *infragrpc.Client) retryfailedmints.BlockchainClient { return c },
		func(db *gorm.DB) processcashbackapproved.TransactionManager {
			return gormtx.NewTransactionManager(db)
		},
	)

	mintInvokes = fx.Invoke(consumer.StartConsumer)

	Mint = fx.Options(
		mintFactories,
		mintDependencies,
		mintInvokes,
	)
)
