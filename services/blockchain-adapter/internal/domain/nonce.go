package domain

import (
	"time"

	"github.com/google/uuid"
)

// WalletNonce tracks nonces for wallet addresses
type WalletNonce struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	WalletAddress string    `gorm:"type:varchar(42);uniqueIndex;not null"`
	CurrentNonce  int64     `gorm:"not null;default:0"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

// TableName specifies the table name for GORM
func (WalletNonce) TableName() string {
	return "wallet_nonces"
}
