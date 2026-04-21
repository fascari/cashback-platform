package domain

import "time"

const (
	StatusPending   Status = "pending"
	StatusProcessed Status = "processed"
	StatusFailed    Status = "failed"
)

type (
	Status string

	Deposit struct {
		ID              int64
		ChainID         string
		TransactionHash string
		WalletAddress   string
		TokenAmount     string
		BlockNumber     int64
		Status          Status
		DetectedAt      time.Time
		ProcessedAt     *time.Time
	}
)
