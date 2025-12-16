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
	// MintRequestStatus represents the status of a mint request
	MintRequestStatus string

	// MintRequest represents a token minting request
	MintRequest struct {
		ID              uuid.UUID         `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
		CashbackID      uuid.UUID         `gorm:"type:uuid;not null;index"`
		UserID          uuid.UUID         `gorm:"type:uuid;not null"`
		WalletAddress   string            `gorm:"type:varchar(42);not null"`
		TokenAmount     string            `gorm:"type:varchar(78);not null"`
		IdempotencyKey  uuid.UUID         `gorm:"type:uuid;uniqueIndex;not null"`
		Status          MintRequestStatus `gorm:"type:varchar(50);not null;default:'pending';index"`
		RetryCount      int               `gorm:"not null;default:0"`
		MaxRetries      int               `gorm:"not null;default:5"`
		TransactionHash string            `gorm:"type:varchar(66)"`
		BlockNumber     int64
		ErrorCode       string `gorm:"type:varchar(100)"`
		ErrorMessage    string `gorm:"type:text"`
		NextRetryAt     *time.Time
		CreatedAt       time.Time `gorm:"autoCreateTime"`
		UpdatedAt       time.Time `gorm:"autoUpdateTime"`
		CompletedAt     *time.Time
	}

	// TokenMintRequestedEvent represents the token.mint.requested domain event
	TokenMintRequestedEvent struct {
		EventID   uuid.UUID `json:"event_id"`
		EventType string    `json:"event_type"`
		Timestamp time.Time `json:"timestamp"`
		Data      struct {
			MintRequestID  uuid.UUID `json:"mint_request_id"`
			CashbackID     uuid.UUID `json:"cashback_id"`
			UserID         uuid.UUID `json:"user_id"`
			WalletAddress  string    `json:"wallet_address"`
			TokenAmount    string    `json:"token_amount"`
			IdempotencyKey uuid.UUID `json:"idempotency_key"`
		} `json:"data"`
	}

	// TokenMintedEvent represents the token.minted domain event
	TokenMintedEvent struct {
		EventID   uuid.UUID `json:"event_id"`
		EventType string    `json:"event_type"`
		Timestamp time.Time `json:"timestamp"`
		Data      struct {
			MintRequestID   uuid.UUID `json:"mint_request_id"`
			CashbackID      uuid.UUID `json:"cashback_id"`
			UserID          uuid.UUID `json:"user_id"`
			WalletAddress   string    `json:"wallet_address"`
			TokenAmount     string    `json:"token_amount"`
			TransactionHash string    `json:"transaction_hash"`
			BlockNumber     int64     `json:"block_number"`
			MintedAt        time.Time `json:"minted_at"`
		} `json:"data"`
	}

	// TokenMintFailedEvent represents the token.mint.failed domain event
	TokenMintFailedEvent struct {
		EventID   uuid.UUID `json:"event_id"`
		EventType string    `json:"event_type"`
		Timestamp time.Time `json:"timestamp"`
		Data      struct {
			MintRequestID uuid.UUID  `json:"mint_request_id"`
			CashbackID    uuid.UUID  `json:"cashback_id"`
			UserID        uuid.UUID  `json:"user_id"`
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

// TableName specifies the table name for GORM
func (MintRequest) TableName() string {
	return "mint_requests"
}

func NewTokenMintRequestedEvent(req *MintRequest) *TokenMintRequestedEvent {
	event := &TokenMintRequestedEvent{
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

func NewTokenMintedEvent(req *MintRequest) *TokenMintedEvent {
	event := &TokenMintedEvent{
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

func NewTokenMintFailedEvent(req *MintRequest) *TokenMintFailedEvent {
	event := &TokenMintFailedEvent{
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
