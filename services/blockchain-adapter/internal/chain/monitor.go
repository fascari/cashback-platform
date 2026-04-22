package chain

import (
	"context"
	"time"
)

type (
	Deposit struct {
		ChainID         ID        `json:"chain_id"`
		TransactionHash string    `json:"transaction_hash"`
		FromAddress     string    `json:"from_address"`
		ToAddress       string    `json:"to_address"`
		TokenAmount     string    `json:"token_amount"`
		BlockNumber     int64     `json:"block_number"`
		DetectedAt      time.Time `json:"detected_at"`
	}

	DepositHandler func(ctx context.Context, d Deposit) error

	DepositMonitor interface {
		Watch(ctx context.Context, handler DepositHandler) error
		Stop()
	}
)
