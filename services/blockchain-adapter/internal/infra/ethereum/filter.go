package ethereum

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/cashback-platform/services/blockchain-adapter/internal/chain"
	"github.com/cashback-platform/services/blockchain-adapter/internal/contracts"
)

// TokenFilter queries ERC-20 Transfer events for a given block range.
// Defined here so it can be mocked in tests without depending on the generated contract struct.
type (
	TokenFilter interface {
		FilterTransfers(ctx context.Context, fromBlock, toBlock uint64) ([]chain.Deposit, error)
	}

	contractFilter struct {
		token *contracts.CashbackToken
	}
)

func NewContractFilter(token *contracts.CashbackToken) TokenFilter {
	return contractFilter{token: token}
}

func (f contractFilter) FilterTransfers(ctx context.Context, fromBlock, toBlock uint64) ([]chain.Deposit, error) {
	iter, err := f.token.FilterTransfer(
		&bind.FilterOpts{Start: fromBlock, End: &toBlock, Context: ctx},
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}
	defer func() { _ = iter.Close() }()

	var deposits []chain.Deposit
	for iter.Next() {
		ev := iter.Event
		deposits = append(deposits, chain.Deposit{
			ChainID:         chain.Ethereum,
			TransactionHash: ev.Raw.TxHash.Hex(),
			FromAddress:     ev.From.Hex(),
			ToAddress:       ev.To.Hex(),
			TokenAmount:     ev.Value.String(),
			BlockNumber:     int64(ev.Raw.BlockNumber),
			DetectedAt:      time.Now().UTC(),
		})
	}
	return deposits, iter.Error()
}
