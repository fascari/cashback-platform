package mintrequest

import (
	"time"

	"github.com/google/uuid"

	"github.com/cashback-platform/services/mint-consumer/internal/domain"
)

type mintRequestModel struct {
	ID              int64                    `gorm:"primaryKey;autoIncrement"`
	CashbackID      int64                    `gorm:"not null;index"`
	UserID          int64                    `gorm:"not null"`
	WalletAddress   string                   `gorm:"type:varchar(42);not null"`
	TokenAmount     string                   `gorm:"type:varchar(78);not null"`
	IdempotencyKey  uuid.UUID                `gorm:"type:uuid;uniqueIndex;not null"`
	Status          domain.MintRequestStatus `gorm:"type:mint.mint_request_status;not null;default:pending;index"`
	RetryCount      int                      `gorm:"not null;default:0"`
	MaxRetries      int                      `gorm:"not null;default:5"`
	TransactionHash string                   `gorm:"type:varchar(66)"`
	BlockNumber     int64
	ErrorCode       string `gorm:"type:varchar(100)"`
	ErrorMessage    string `gorm:"type:text"`
	NextRetryAt     *time.Time
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`
	CompletedAt     *time.Time
}

func (mintRequestModel) TableName() string {
	return "mint.mint_requests"
}

func (m mintRequestModel) toDomain() domain.MintRequest {
	return domain.MintRequest{
		ID:              m.ID,
		CashbackID:      m.CashbackID,
		UserID:          m.UserID,
		WalletAddress:   m.WalletAddress,
		TokenAmount:     m.TokenAmount,
		IdempotencyKey:  m.IdempotencyKey,
		Status:          m.Status,
		RetryCount:      m.RetryCount,
		MaxRetries:      m.MaxRetries,
		TransactionHash: m.TransactionHash,
		BlockNumber:     m.BlockNumber,
		ErrorCode:       m.ErrorCode,
		ErrorMessage:    m.ErrorMessage,
		NextRetryAt:     m.NextRetryAt,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
		CompletedAt:     m.CompletedAt,
	}
}

func fromDomain(r domain.MintRequest) mintRequestModel {
	return mintRequestModel{
		ID:              r.ID,
		CashbackID:      r.CashbackID,
		UserID:          r.UserID,
		WalletAddress:   r.WalletAddress,
		TokenAmount:     r.TokenAmount,
		IdempotencyKey:  r.IdempotencyKey,
		Status:          r.Status,
		RetryCount:      r.RetryCount,
		MaxRetries:      r.MaxRetries,
		TransactionHash: r.TransactionHash,
		BlockNumber:     r.BlockNumber,
		ErrorCode:       r.ErrorCode,
		ErrorMessage:    r.ErrorMessage,
		NextRetryAt:     r.NextRetryAt,
		CompletedAt:     r.CompletedAt,
	}
}
