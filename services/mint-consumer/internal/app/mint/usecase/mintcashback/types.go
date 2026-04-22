package mintcashback

import "github.com/google/uuid"

type Input struct {
	EventID       uuid.UUID
	CashbackID    int64
	UserID        int64
	WalletAddress string
	PurchaseID    string
	TokenAmount   string
	ChainID       string
}
