package domain

import "time"

// WalletNonce tracks nonces for wallet addresses
type WalletNonce struct {
	ID            int64     `gorm:"primaryKey;autoIncrement"`
	WalletAddress string    `gorm:"type:varchar(42);uniqueIndex;not null"`
	CurrentNonce  int64     `gorm:"not null;default:0"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

// TableName specifies the schema-qualified table name for GORM.
func (WalletNonce) TableName() string {
	return "blockchain.wallet_nonces"
}
