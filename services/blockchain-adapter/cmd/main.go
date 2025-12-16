package main

import (
	"github.com/cashback-platform/services/blockchain-adapter/internal/config"
	grpcserver "github.com/cashback-platform/services/blockchain-adapter/internal/grpc"
	"github.com/cashback-platform/services/blockchain-adapter/internal/infra/database"
	repoNonce "github.com/cashback-platform/services/blockchain-adapter/internal/repository/nonce"
	repoTransaction "github.com/cashback-platform/services/blockchain-adapter/internal/repository/transaction"
	usecaseToken "github.com/cashback-platform/services/blockchain-adapter/internal/usecase"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		// Configuration
		fx.Provide(config.NewConfig),

		// Infrastructure
		fx.Provide(database.NewPostgresDB),

		// Repositories
		fx.Provide(repoTransaction.NewRepository),
		fx.Provide(repoNonce.NewRepository),

		// Usecases
		fx.Provide(usecaseToken.NewTokenUsecase),

		// gRPC Server
		fx.Provide(grpcserver.NewTokenServer),

		// Start server
		fx.Invoke(grpcserver.StartServer),
	).Run()
}
