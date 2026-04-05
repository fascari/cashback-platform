package modules

import (
	"context"

	"github.com/cashback-platform/kit/gormtx"
	"github.com/cashback-platform/services/mint-consumer/internal/app/mint/repository"
	"github.com/cashback-platform/services/mint-consumer/internal/consumer"
	infragrpc "github.com/cashback-platform/services/mint-consumer/internal/infra/grpc"
	"github.com/cashback-platform/services/mint-consumer/internal/usecase"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

var (
	mintFactories = fx.Provide(
		repository.New,
		usecase.NewMint,
		consumer.NewCashback,
	)

	mintDependencies = fx.Provide(
		func(r repository.Repository) usecase.MintRequestRepository { return r },
		func(r repository.Repository) usecase.ProcessedEventRepository { return r },
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


