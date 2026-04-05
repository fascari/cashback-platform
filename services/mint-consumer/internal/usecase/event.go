package usecase

import (
	"time"

	"github.com/google/uuid"

	"github.com/cashback-platform/services/mint-consumer/internal/domain"
)

type (
	CashbackApprovedEvent struct {
		EventID       string `json:"event_id"`
		CashbackID    string `json:"cashback_id"`
		UserID        string `json:"user_id"`
		WalletAddress string `json:"wallet_address"`
		PurchaseID    string `json:"purchase_id"`
		TokenAmount   string `json:"token_amount"`
		ChainID       string `json:"chain_id"`
	}

	TokenMintRequestedEvent struct {
		EventID   uuid.UUID `json:"event_id"`
		EventType string    `json:"event_type"`
		Timestamp time.Time `json:"timestamp"`
		Data      struct {
			MintRequestID  int64     `json:"mint_request_id"`
			CashbackID     int64     `json:"cashback_id"`
			UserID         int64     `json:"user_id"`
			WalletAddress  string    `json:"wallet_address"`
			TokenAmount    string    `json:"token_amount"`
			IdempotencyKey uuid.UUID `json:"idempotency_key"`
		} `json:"data"`
	}

	TokenMintedEvent struct {
		EventID   uuid.UUID `json:"event_id"`
		EventType string    `json:"event_type"`
		Timestamp time.Time `json:"timestamp"`
		Data      struct {
			MintRequestID   int64     `json:"mint_request_id"`
			CashbackID      int64     `json:"cashback_id"`
			UserID          int64     `json:"user_id"`
			WalletAddress   string    `json:"wallet_address"`
			TokenAmount     string    `json:"token_amount"`
			TransactionHash string    `json:"transaction_hash"`
			BlockNumber     int64     `json:"block_number"`
			MintedAt        time.Time `json:"minted_at"`
		} `json:"data"`
	}

	TokenMintFailedEvent struct {
		EventID   uuid.UUID `json:"event_id"`
		EventType string    `json:"event_type"`
		Timestamp time.Time `json:"timestamp"`
		Data      struct {
			MintRequestID int64      `json:"mint_request_id"`
			CashbackID    int64      `json:"cashback_id"`
			UserID        int64      `json:"user_id"`
			WalletAddress string     `json:"wallet_address"`
			TokenAmount   string     `json:"token_amount"`
			ErrorCode     string     `json:"error_code"`
			ErrorMessage  string     `json:"error_message"`
			RetryCount    int        `json:"retry_count"`
			MaxRetries    int        `json:"max_retries"`
			NextRetryAt   *time.Time `json:"next_retry_at,omitempty"`
		} `json:"data"`
	}
)

func NewTokenMintRequestedEvent(req domain.MintRequest) TokenMintRequestedEvent {
	event := TokenMintRequestedEvent{
		EventID:   uuid.New(),
		EventType: "token.mint.requested",
		Timestamp: time.Now().UTC(),
	}
	event.Data.MintRequestID = req.ID
	event.Data.CashbackID = req.CashbackID
	event.Data.UserID = req.UserID
	event.Data.WalletAddress = req.WalletAddress
	event.Data.TokenAmount = req.TokenAmount
	event.Data.IdempotencyKey = req.IdempotencyKey
	return event
}

func NewTokenMintedEvent(req domain.MintRequest) TokenMintedEvent {
	event := TokenMintedEvent{
		EventID:   uuid.New(),
		EventType: "token.minted",
		Timestamp: time.Now().UTC(),
	}
	event.Data.MintRequestID = req.ID
	event.Data.CashbackID = req.CashbackID
	event.Data.UserID = req.UserID
	event.Data.WalletAddress = req.WalletAddress
	event.Data.TokenAmount = req.TokenAmount
	event.Data.TransactionHash = req.TransactionHash
	event.Data.BlockNumber = req.BlockNumber
	event.Data.MintedAt = time.Now().UTC()
	return event
}

func NewTokenMintFailedEvent(req domain.MintRequest) TokenMintFailedEvent {
	event := TokenMintFailedEvent{
		EventID:   uuid.New(),
		EventType: "token.mint.failed",
		Timestamp: time.Now().UTC(),
	}
	event.Data.MintRequestID = req.ID
	event.Data.CashbackID = req.CashbackID
	event.Data.UserID = req.UserID
	event.Data.WalletAddress = req.WalletAddress
	event.Data.TokenAmount = req.TokenAmount
	event.Data.ErrorCode = req.ErrorCode
	event.Data.ErrorMessage = req.ErrorMessage
	event.Data.RetryCount = req.RetryCount
	event.Data.MaxRetries = req.MaxRetries
	event.Data.NextRetryAt = req.NextRetryAt
	return event
}
