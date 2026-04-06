package modules

import (
	"github.com/cashback-platform/kit/gormtx"
	"github.com/cashback-platform/services/mint-consumer/internal/app/mint/repository"
	"github.com/cashback-platform/services/mint-consumer/internal/app/mint/usecase/mintcashback"
	"github.com/cashback-platform/services/mint-consumer/internal/app/mint/usecase/retrymints"
	"github.com/cashback-platform/services/mint-consumer/internal/consumer"
	"github.com/cashback-platform/services/mint-consumer/internal/consumer/cashbackapproved"
	infragrpc "github.com/cashback-platform/services/mint-consumer/internal/infra/grpc"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var (
	mintFactories = fx.Provide(
		repository.New,
		mintcashback.NewUseCase,
		retrymints.NewUseCase,
		cashbackapproved.New,
		consumer.NewCashback,
	)

	mintDependencies = fx.Provide(
		func(r repository.Repository) mintcashback.Repository { return r },
		func(r repository.Repository) retrymints.Repository { return r },
		func(c *infragrpc.Client) mintcashback.BlockchainClient { return c },
		func(c *infragrpc.Client) retrymints.BlockchainClient { return c },
		func(db *gorm.DB) mintcashback.TransactionManager {
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
