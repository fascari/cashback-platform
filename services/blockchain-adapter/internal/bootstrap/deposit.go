package bootstrap

import (
	"context"
	"encoding/json"
	"fmt"

	"go.uber.org/fx"

	ethereumpkg "github.com/cashback-platform/kit/ethereum"
	kitev "github.com/cashback-platform/kit/events"
	"github.com/cashback-platform/kit/logger"
	"github.com/cashback-platform/services/blockchain-adapter/internal/app/deposit/domain"
	depositrepo "github.com/cashback-platform/services/blockchain-adapter/internal/app/deposit/repository"
	"github.com/cashback-platform/services/blockchain-adapter/internal/chain"
	"github.com/cashback-platform/services/blockchain-adapter/internal/contracts"
	ethereuminfra "github.com/cashback-platform/services/blockchain-adapter/internal/infra/ethereum"
	infranats "github.com/cashback-platform/services/blockchain-adapter/internal/infra/nats"
)

var Deposit = fx.Module("deposit",
	fx.Provide(
		depositrepo.New,
		newDepositMonitor,
	),
	fx.Invoke(startDepositMonitor),
)

type (
	depositMonitorParams struct {
		fx.In

		Token  *contracts.CashbackToken
		Client *ethereumpkg.Client
		Repo   depositrepo.Repository
	}

	depositLifecycleParams struct {
		fx.In

		LC      fx.Lifecycle
		Monitor *ethereuminfra.DepositMonitor
		Repo    depositrepo.Repository
		NATS    *infranats.NATSClient
	}
)

func newDepositMonitor(p depositMonitorParams) (*ethereuminfra.DepositMonitor, error) {
	ctx := context.Background()

	startBlock, err := p.Repo.MaxBlockNumber(ctx, string(chain.Ethereum))
	if err != nil {
		return nil, fmt.Errorf("query max block number: %w", err)
	}

	if startBlock > 0 {
		startBlock++
	}

	filter := ethereuminfra.NewContractFilter(p.Token)
	return ethereuminfra.NewDepositMonitor(filter, p.Client.Inner(), uint64(startBlock)), nil //nolint:gosec
}

func startDepositMonitor(p depositLifecycleParams) {
	handler := newDepositHandler(p.Repo, p.NATS)

	// monitorCtx must outlive the FX OnStart context, which expires after startup.
	monitorCtx, cancel := context.WithCancel(context.Background())

	p.LC.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() {
				if err := p.Monitor.Watch(monitorCtx, handler); err != nil {
					logger.Error("deposit monitor stopped", "error", err)
				}
			}()
			return nil
		},
		OnStop: func(_ context.Context) error {
			cancel()
			p.Monitor.Stop()
			return nil
		},
	})
}

func newDepositHandler(repo depositrepo.Repository, nats *infranats.NATSClient) chain.DepositHandler {
	return func(ctx context.Context, d chain.Deposit) error {
		if err := repo.Save(ctx, domain.Deposit{
			ChainID:         string(d.ChainID),
			TransactionHash: d.TransactionHash,
			WalletAddress:   d.ToAddress,
			TokenAmount:     d.TokenAmount,
			BlockNumber:     d.BlockNumber,
			Status:          domain.StatusPending,
			DetectedAt:      d.DetectedAt,
		}); err != nil {
			return fmt.Errorf("save deposit %s: %w", d.TransactionHash, err)
		}

		payload, err := json.Marshal(d)
		if err != nil {
			return fmt.Errorf("marshal deposit %s: %w", d.TransactionHash, err)
		}

		if err := nats.Publish(kitev.DepositDetected, payload); err != nil {
			return fmt.Errorf("publish deposit %s: %w", d.TransactionHash, err)
		}

		logger.Info("deposit detected and published",
			"tx_hash", d.TransactionHash,
			"chain_id", d.ChainID,
			"block_number", d.BlockNumber,
		)
		return nil
	}
}
