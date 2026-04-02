package main

import (
	"github.com/cashback-platform/services/cashback-service-api/cmd/api/modules"
	"github.com/cashback-platform/services/cashback-service-api/internal/bootstrap"
	"github.com/cashback-platform/services/cashback-service-api/internal/config"
	"github.com/cashback-platform/services/cashback-service-api/internal/infra/grpc"
	"github.com/cashback-platform/services/cashback-service-api/internal/infra/messaging/outbox"
	"github.com/cashback-platform/services/cashback-service-api/internal/infra/messaging/outbox/relay"
	outboxrepo "github.com/cashback-platform/services/cashback-service-api/internal/infra/messaging/outbox/repository"
	"github.com/cashback-platform/services/cashback-service-api/internal/infra/nats"
	"github.com/cashback-platform/services/cashback-service-api/pkg/logger"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

func main() {
	logger.Init()

	app := fx.New(
		bootstrap.Logger(),
		bootstrap.Config,
		bootstrap.Database,
		fx.Provide(nats.NewNATSClient),
		fx.Provide(grpc.NewBlockchainAdapterClient),
		bootstrap.Router,
		bootstrap.Server,
		fx.Provide(config.LoadOutbox),
		fx.Provide(func(db *gorm.DB, cfg config.OutboxConfig) outboxrepo.Repository {
			return outboxrepo.New(db, cfg.MaxRetries)
		}),
		fx.Provide(func(repo outboxrepo.Repository) outbox.Publisher {
			return outbox.New(repo)
		}),
		fx.Provide(func(repo outboxrepo.Repository, natsClient *nats.NATSClient, cfg config.OutboxConfig) relay.Relay {
			return relay.New(repo, natsClient, cfg.PollInterval)
		}),
		fx.Invoke(relay.Start),
		modules.User,
		modules.Purchase,
		modules.Cashback,
	)

	app.Run()
}
