package main

import (
	"github.com/cashback-platform/services/mint-consumer/internal/config"
	"github.com/cashback-platform/services/mint-consumer/internal/consumer"
	"github.com/cashback-platform/services/mint-consumer/internal/infra/database"
	"github.com/cashback-platform/services/mint-consumer/internal/infra/grpc"
	"github.com/cashback-platform/services/mint-consumer/internal/infra/nats"
	repoMintRequest "github.com/cashback-platform/services/mint-consumer/internal/repository/mintrequest"
	repoProcessedEvent "github.com/cashback-platform/services/mint-consumer/internal/repository/processedevent"
	"github.com/cashback-platform/services/mint-consumer/internal/usecase"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		// Configuration
		fx.Provide(config.NewConfig),

		// Infrastructure
		fx.Provide(database.NewPostgresDB),
		fx.Provide(nats.NewNATSClient),
		fx.Provide(grpc.NewBlockchainAdapterClient),

		// Repositories
		fx.Provide(repoMintRequest.NewRepository),
		fx.Provide(repoProcessedEvent.NewRepository),

		// Usecases
		fx.Provide(usecase.NewMintUsecase),

		// Consumer
		fx.Provide(consumer.NewCashbackConsumer),

		// Start consumer
		fx.Invoke(consumer.StartConsumer),
	).Run()
}
