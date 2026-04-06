package domain

import (
	"time"

	"github.com/cashback-platform/kit/apperror"
	"github.com/google/uuid"
)

const (
	ErrCodeMintRequestNotFound  = "error_mint_request_not_found"
	ErrCodeMintRequestDuplicate = "error_mint_request_duplicate"

	MintRequestStatusPending    MintRequestStatus = "pending"
	MintRequestStatusProcessing MintRequestStatus = "processing"
	MintRequestStatusCompleted  MintRequestStatus = "completed"
	MintRequestStatusFailed     MintRequestStatus = "failed"
)

var (
	ErrMintRequestNotFound  = apperror.New(ErrCodeMintRequestNotFound, "mint request not found")
	ErrMintRequestDuplicate = apperror.New(ErrCodeMintRequestDuplicate, "mint request already exists")
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
