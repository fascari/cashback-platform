package nonce

import (
	"time"
)

type walletNonceModel struct {
	ID            int64     `gorm:"primaryKey;autoIncrement"`
	WalletAddress string    `gorm:"type:varchar(42);uniqueIndex;not null"`
	CurrentNonce  int64     `gorm:"not null;default:0"`
	FenceToken    int64     `gorm:"not null;default:0"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

func (walletNonceModel) TableName() string {
	return "blockchain.wallet_nonces"
}
