package ethereum

//go:generate mockery --all --case=snake --disable-version-string --with-expecter

import (
	"context"
	"sync"
	"time"

	"github.com/cashback-platform/kit/logger"
	"github.com/cashback-platform/services/blockchain-adapter/internal/chain"
)

const defaultPollInterval = 12 * time.Second

type (
	BlockReader interface {
		BlockNumber(ctx context.Context) (uint64, error)
	}

	// DepositMonitor polls an Ethereum node for ERC-20 Transfer events
	// and calls the provided DepositHandler for every new event detected.
	DepositMonitor struct {
		filter   TokenFilter
		blocks   BlockReader
		interval time.Duration
		start    uint64
		quit     chan struct{}
		once     sync.Once
	}
)

// NewDepositMonitor creates a monitor that starts polling from startBlock.
// Pass 0 to let the monitor seed from (latest - 1000) on first poll.
func NewDepositMonitor(filter TokenFilter, blocks BlockReader, startBlock uint64, opts ...func(*DepositMonitor)) *DepositMonitor {
	m := new(DepositMonitor{
		filter:   filter,
		blocks:   blocks,
		interval: defaultPollInterval,
		start:    startBlock,
		quit:     make(chan struct{}),
	})
	for _, o := range opts {
		o(m)
	}
	return m
}

func WithPollInterval(d time.Duration) func(*DepositMonitor) {
	return func(m *DepositMonitor) { m.interval = d }
}

// Watch polls for Transfer events and calls handler for each detected deposit.
// Blocks until ctx is cancelled or Stop is called.
func (m *DepositMonitor) Watch(ctx context.Context, handler chain.DepositHandler) error {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	current := m.start

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-m.quit:
			return nil
		case <-ticker.C:
			current = m.tick(ctx, current, handler)
		}
	}
}

func (m *DepositMonitor) tick(ctx context.Context, current uint64, handler chain.DepositHandler) uint64 {
	latest, err := m.blocks.BlockNumber(ctx)
	if err != nil {
		logger.Warn("deposit monitor: get block number failed", "error", err)
		return current
	}

	if current == 0 {
		current = seedBlock(latest)
	}

	// Chain reset: Anvil restarted while the DB still has a higher block number.
	// A difference of exactly 1 is normal "caught up" state, not a reset.
	if current > 0 && latest+1 < current {
		current = seedBlock(latest)
	}

	// No new blocks since last tick.
	if latest < current {
		return current
	}

	deposits, err := m.filter.FilterTransfers(ctx, current, latest)
	if err != nil {
		logger.Warn("deposit monitor: filter transfers failed",
			"from_block", current,
			"to_block", latest,
			"error", err,
		)
		return current
	}

	dispatchDeposits(ctx, handler, deposits)
	return latest + 1
}

func dispatchDeposits(ctx context.Context, handler chain.DepositHandler, deposits []chain.Deposit) {
	for _, d := range deposits {
		if err := handler(ctx, d); err != nil {
			logger.Error("deposit monitor: handler error",
				"tx_hash", d.TransactionHash,
				"error", err,
			)
		}
	}
}

func seedBlock(latest uint64) uint64 {
	if latest > 1000 {
		return latest - 1000
	}
	return 1
}

func (m *DepositMonitor) Stop() {
	m.once.Do(func() {
		close(m.quit)
	})
}
