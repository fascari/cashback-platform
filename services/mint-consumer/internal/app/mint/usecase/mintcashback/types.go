package mintcashback

import "github.com/google/uuid"

// Input carries the event data needed to process a cashback mint.
type Input struct {
	EventID       uuid.UUID
	CashbackID    int64
	UserID        int64
	WalletAddress string
	PurchaseID    string
	TokenAmount   string
	ChainID       string
}
