package domain

import (
	"time"

	"github.com/google/uuid"
)

const (
	MintRequestStatusPending    MintRequestStatus = "pending"
	MintRequestStatusProcessing MintRequestStatus = "processing"
	MintRequestStatusCompleted  MintRequestStatus = "completed"
	MintRequestStatusFailed     MintRequestStatus = "failed"
)

type (
	MintRequestStatus string

	MintRequest struct {
		ID              int64
		CashbackID      int64
		UserID          int64
		WalletAddress   string
		TokenAmount     string
		IdempotencyKey  uuid.UUID
		Status          MintRequestStatus
		RetryCount      int
		MaxRetries      int
		TransactionHash string
		BlockNumber     int64
		ErrorCode       string
		ErrorMessage    string
		NextRetryAt     *time.Time
		CreatedAt       time.Time
		UpdatedAt       time.Time
		CompletedAt     *time.Time
	}
)
