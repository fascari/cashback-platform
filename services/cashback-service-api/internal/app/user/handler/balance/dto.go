package balance

import (
	"strconv"

	balanceuc "github.com/cashback-platform/services/cashback-service-api/internal/app/user/usecase/balance"
)

type OutputPayload struct {
	UserID        string `json:"user_id"`
	WalletAddress string `json:"wallet_address"`
	Balance       string `json:"balance"`
	BalanceTokens string `json:"balance_tokens"`
	BlockNumber   int64  `json:"block_number"`
}

func ToOutputPayload(output balanceuc.Output) OutputPayload {
	return OutputPayload{
		UserID:        strconv.FormatInt(output.UserID, 10),
		WalletAddress: output.WalletAddress,
		Balance:       output.Balance,
		BalanceTokens: output.BalanceTokens,
		BlockNumber:   output.BlockNumber,
	}
}
