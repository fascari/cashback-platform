package modules

import (
	"context"

	"github.com/cashback-platform/kit/gormtx"
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
		func(client *infragrpc.Client) usecase.BlockchainClient {
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

// blockchainGRPCClient adapts infragrpc.Client to usecase.BlockchainClient.
// Removed in Phase 6 once use cases reference infragrpc types directly.
type blockchainGRPCClient struct {
	client *infragrpc.Client
}

func (c blockchainGRPCClient) MintToken(ctx context.Context, req usecase.MintTokenRequest) (usecase.MintResult, error) {
	result, err := c.client.MintToken(ctx, infragrpc.MintTokenRequest{
		IdempotencyKey: req.IdempotencyKey,
		WalletAddress:  req.WalletAddress,
		TokenAmount:    req.TokenAmount,
		ChainID:        req.ChainID,
	})
	if err != nil {
		return usecase.MintResult{}, err
	}
	return usecase.MintResult{
		TransactionHash: result.TransactionHash,
		BlockNumber:     result.BlockNumber,
		ErrorCode:       result.ErrorCode,
		ErrorMessage:    result.ErrorMessage,
		Retryable:       result.Retryable,
	}, nil
}

