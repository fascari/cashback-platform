package balance

type (
	TokenBalance struct {
		WalletAddress string
		Amount        string
		BlockNumber   int64
	}

	Output struct {
		UserID        int64
		WalletAddress string
		Balance       string
		BalanceTokens string
		BlockNumber   int64
	}
)
