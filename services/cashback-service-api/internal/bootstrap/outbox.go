package bootstrap

import (
	"github.com/cashback-platform/services/cashback-service-api/internal/config"
	"github.com/cashback-platform/services/cashback-service-api/internal/infra/messaging/outbox"
	"github.com/cashback-platform/services/cashback-service-api/internal/infra/messaging/outbox/relay"
	outboxrepo "github.com/cashback-platform/services/cashback-service-api/internal/infra/messaging/outbox/repository"
	"github.com/cashback-platform/services/cashback-service-api/internal/infra/nats"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

var Outbox = fx.Module("outbox",
	fx.Provide(
		newOutboxRepository,
		newOutboxPublisher,
		newOutboxRelay,
	),
	fx.Invoke(relay.Start),
)

func newOutboxRepository(db *gorm.DB, cfg config.OutboxConfig) outboxrepo.Repository {
	return outboxrepo.New(db, cfg.MaxRetries)
}

func newOutboxPublisher(repo outboxrepo.Repository) outbox.Publisher {
	return outbox.New(repo)
}

func newOutboxRelay(repo outboxrepo.Repository, natsClient *nats.NATSClient, cfg config.OutboxConfig) relay.Relay {
	return relay.New(repo, natsClient, cfg.PollInterval)
}
