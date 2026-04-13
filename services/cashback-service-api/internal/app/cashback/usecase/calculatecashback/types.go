package calculatecashback

// CashbackApprovedEvent is the payload published to NATS when a cashback is approved.
type CashbackApprovedEvent struct {
	EventID       string `json:"event_id"`
	CashbackID    string `json:"cashback_id"`
	UserID        string `json:"user_id"`
	WalletAddress string `json:"wallet_address"`
	PurchaseID    string `json:"purchase_id"`
	TokenAmount   string `json:"token_amount"`
	ChainID       string `json:"chain_id"`
}
