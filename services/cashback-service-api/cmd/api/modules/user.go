package modules

import (
	"context"

	"github.com/cashback-platform/services/cashback-service-api/internal/app/user/handler/balance"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/user/handler/createuser"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/user/handler/finduser"
	userrepo "github.com/cashback-platform/services/cashback-service-api/internal/app/user/repository"
	balanceuc "github.com/cashback-platform/services/cashback-service-api/internal/app/user/usecase/balance"
	createuseruc "github.com/cashback-platform/services/cashback-service-api/internal/app/user/usecase/createuser"
	finduseruc "github.com/cashback-platform/services/cashback-service-api/internal/app/user/usecase/finduser"
	grpcclient "github.com/cashback-platform/services/cashback-service-api/internal/infra/grpc"

	"go.uber.org/fx"
)

var (
	userFactories = fx.Provide(
		userrepo.New,
		createuseruc.New,
		finduseruc.New,
		balanceuc.New,
		createuser.NewHandler,
		finduser.NewHandler,
		balance.NewHandler,
	)

	userDependencies = fx.Provide(
		func(repo userrepo.Repository) createuseruc.Repository {
			return repo
		},
		func(repo userrepo.Repository) finduseruc.Repository {
			return repo
		},
		func(repo userrepo.Repository) balanceuc.UserRepository {
			return repo
		},
		func(client *grpcclient.BlockchainAdapterClient) balanceuc.BlockchainClient {
			return blockchainClientAdapter{client}
		},
	)

	userInvokes = fx.Invoke(
		func(params RouterParams, h createuser.Handler) {
			createuser.RegisterEndpoint(params.APIRouter, h)
		},
		func(params RouterParams, h finduser.Handler) {
			finduser.RegisterEndpoint(params.APIRouter, h)
		},
		func(params RouterParams, h balance.Handler) {
			balance.RegisterEndpoint(params.APIRouter, h)
		},
	)

	User = fx.Options(
		userFactories,
		userDependencies,
		userInvokes,
	)
)

// blockchainClientAdapter adapts grpcclient.BlockchainAdapterClient to balanceuc.BlockchainClient.
type blockchainClientAdapter struct {
	client *grpcclient.BlockchainAdapterClient
}

func (a blockchainClientAdapter) Balance(ctx context.Context, walletAddress string) (balanceuc.TokenBalance, error) {
	resp, err := a.client.Balance(ctx, walletAddress)
	if err != nil {
		return balanceuc.TokenBalance{}, err
	}
	return balanceuc.TokenBalance{
		WalletAddress: resp.WalletAddress,
		Amount:        resp.Balance,
		BlockNumber:   resp.BlockNumber,
	}, nil
}
