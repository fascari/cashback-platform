package calculatecashback

// Purchase is a local projection of the purchase data needed for cashback calculation.
// User is a local projection of the user data needed for cashback calculation.
// CashbackApprovedEvent is the payload published to NATS when a cashback is approved.
type (
	Purchase struct {
		ID     int64
		UserID int64
		Amount float64
	}

	User struct {
		WalletAddress string
	}

	CashbackApprovedEvent struct {
		EventID       string `json:"event_id"`
		CashbackID    string `json:"cashback_id"`
		UserID        string `json:"user_id"`
		WalletAddress string `json:"wallet_address"`
		PurchaseID    string `json:"purchase_id"`
		TokenAmount   string `json:"token_amount"`
		ChainID       string `json:"chain_id"`
	}
)
