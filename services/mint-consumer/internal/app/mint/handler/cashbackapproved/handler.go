package cashbackapproved

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cashback-platform/services/mint-consumer/internal/app/mint/usecase/mintcashback"
	natsgo "github.com/nats-io/nats.go"
)

type Handler struct {
	useCase mintcashback.UseCase
}

func New(useCase mintcashback.UseCase) Handler {
	return Handler{useCase: useCase}
}

func (h Handler) Handle(ctx context.Context, msg *natsgo.Msg) error {
	p := new(cashbackApprovedPayload)
	if err := json.Unmarshal(msg.Data, p); err != nil {
		return fmt.Errorf("unmarshal cashback approved payload: %w", err)
	}

	input, err := p.toDomain()
	if err != nil {
		return fmt.Errorf("parse cashback approved payload: %w", err)
	}

	return h.useCase.Execute(ctx, input)
}
