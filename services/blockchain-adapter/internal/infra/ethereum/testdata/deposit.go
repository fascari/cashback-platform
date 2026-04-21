package testdata

import (
	"time"

	"github.com/cashback-platform/services/blockchain-adapter/internal/chain"
)

func DetectedDeposit() chain.Deposit {
	return chain.Deposit{
		ChainID:         chain.Ethereum,
		TransactionHash: "0xabc",
		FromAddress:     "0xfrom",
		ToAddress:       "0xto",
		TokenAmount:     "1000000000000000000",
		BlockNumber:     100,
		DetectedAt:      time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	}
}
