package outbox

import (
	"context"
	"encoding/json"

	"github.com/cashback-platform/services/cashback-service-api/internal/infra/messaging/outbox/repository"
)

type Publisher struct {
	repo repository.Repository
}

func New(repo repository.Repository) Publisher {
	return Publisher{repo: repo}
}

func (p Publisher) Publish(ctx context.Context, eventType string, aggregateID int64, payload any) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return p.repo.Create(ctx, eventType, "cashback", aggregateID, payloadBytes)
}
