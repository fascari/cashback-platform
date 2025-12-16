package main

import (
	"github.com/cashback-platform/services/cashback-service-api/cmd/api/modules"
	"github.com/cashback-platform/services/cashback-service-api/internal/bootstrap"
	"github.com/cashback-platform/services/cashback-service-api/internal/infra/grpc"
	"github.com/cashback-platform/services/cashback-service-api/internal/infra/messaging"
	"github.com/cashback-platform/services/cashback-service-api/internal/infra/messaging/outbox"
	outboxrepo "github.com/cashback-platform/services/cashback-service-api/internal/infra/messaging/outbox/repository"
	"github.com/cashback-platform/services/cashback-service-api/internal/infra/nats"
	"github.com/cashback-platform/services/cashback-service-api/pkg/logger"

	"go.uber.org/fx"
)

func main() {
	logger.Init()

	app := fx.New(
		bootstrap.Logger(),
		// Infrastructure
		bootstrap.Config,
		bootstrap.Database,
		fx.Provide(nats.NewNATSClient),
		fx.Provide(grpc.NewBlockchainAdapterClient),
		bootstrap.Router,
		bootstrap.Server,
		// Messaging (Outbox Pattern)
		fx.Provide(outboxrepo.New),
		fx.Provide(outbox.NewOutboxPublisher),
		fx.Provide(func(op *outbox.OutboxPublisher) messaging.EventPublisher {
			return op
		}),
		fx.Invoke(outbox.StartOutboxPublisher),
		// Business Modules
		modules.User,
		modules.Purchase,
		modules.Cashback,
	)

	app.Run()
}
