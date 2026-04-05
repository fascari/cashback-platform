package modules

import (
	"context"
	"fmt"

	"github.com/cashback-platform/kit/gormtx"
	tokenpb "github.com/cashback-platform/proto/token"
	"github.com/cashback-platform/services/mint-consumer/internal/consumer"
	infragrpc "github.com/cashback-platform/services/mint-consumer/internal/infra/grpc"
	repoMintRequest "github.com/cashback-platform/services/mint-consumer/internal/repository/mintrequest"
	repoProcessedEvent "github.com/cashback-platform/services/mint-consumer/internal/repository/processedevent"
	"github.com/cashback-platform/services/mint-consumer/internal/usecase"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

var (
	mintFactories = fx.Provide(
		repoMintRequest.NewRepository,
		repoProcessedEvent.NewRepository,
		usecase.NewMint,
		consumer.NewCashback,
	)

	mintDependencies = fx.Provide(
		func(client *infragrpc.BlockchainAdapterClient) usecase.BlockchainClient {
			return blockchainGRPCClient{client: client}
		},
		func(db *gorm.DB) usecase.TransactionManager {
			return gormtx.NewTransactionManager(db)
		},
	)

	mintInvokes = fx.Invoke(consumer.StartConsumer)

	Mint = fx.Options(
		mintFactories,
		mintDependencies,
		mintInvokes,
	)
)

type blockchainGRPCClient struct {
	client *infragrpc.BlockchainAdapterClient
}

func (c blockchainGRPCClient) MintToken(ctx context.Context, req usecase.MintTokenRequest) (usecase.MintResult, error) {
	resp, err := c.client.MintToken(ctx, req.IdempotencyKey, req.WalletAddress, req.TokenAmount)
	if err != nil {
		return usecase.MintResult{}, fmt.Errorf("mint token grpc: %w", err)
	}

	result := usecase.MintResult{
		TransactionHash: resp.GetTransactionHash(),
		BlockNumber:     resp.GetBlockNumber(),
	}

	if e := resp.GetError(); e != nil {
		result.ErrorCode = e.GetCode()
		result.ErrorMessage = e.GetMessage()
		result.Retryable = e.GetRetryable()
	}

	if resp.GetStatus() == tokenpb.MintStatus_MINT_STATUS_FAILED && result.ErrorCode == "" {
		result.ErrorCode = "mint_failed"
		result.ErrorMessage = "mint operation failed without details"
	}

	return result, nil
}
