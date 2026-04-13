package modules

import (
	"go.uber.org/fx"
	"google.golang.org/grpc"

	ethereumpkg "github.com/cashback-platform/kit/ethereum"
	"github.com/cashback-platform/services/blockchain-adapter/internal/app/token/handler"
	repoNonce "github.com/cashback-platform/services/blockchain-adapter/internal/app/token/repository/nonce"
	repoTransaction "github.com/cashback-platform/services/blockchain-adapter/internal/app/token/repository/transaction"
	usecaseToken "github.com/cashback-platform/services/blockchain-adapter/internal/app/token/usecase"
	"github.com/cashback-platform/services/blockchain-adapter/internal/contracts"
)

var (
	Token = fx.Options(
		tokenFactories,
		tokenDependencies,
		tokenInvokes,
	)

	tokenFactories = fx.Provide(
		repoTransaction.NewRepository,
		repoNonce.NewRepository,
		usecaseToken.NewToken,
		handler.NewHandler,
	)

	tokenDependencies = fx.Provide(
		func(r repoNonce.Repository) usecaseToken.NonceRepository { return r },
		func(r repoTransaction.Repository) usecaseToken.TransactionRepository { return r },
		func(c *ethereumpkg.Client) usecaseToken.EthereumClient { return c },
		func(t *contracts.CashbackToken) usecaseToken.TokenContract { return t },
	)

	tokenInvokes = fx.Invoke(
		func(s *grpc.Server, h handler.Handler) {
			handler.RegisterServer(s, h)
		},
	)
)
