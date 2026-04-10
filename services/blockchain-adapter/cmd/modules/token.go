package modules

import (
	"go.uber.org/fx"

	ethereumpkg "github.com/cashback-platform/kit/ethereum"
	repoNonce "github.com/cashback-platform/services/blockchain-adapter/internal/app/token/repository/nonce"
	repoTransaction "github.com/cashback-platform/services/blockchain-adapter/internal/app/token/repository/transaction"
	usecaseToken "github.com/cashback-platform/services/blockchain-adapter/internal/app/token/usecase"
	"github.com/cashback-platform/services/blockchain-adapter/internal/contracts"
)

var (
	Token = fx.Options(
		tokenFactories,
		tokenDependencies,
	)

	tokenFactories = fx.Provide(
		repoTransaction.NewRepository,
		repoNonce.NewRepository,
		usecaseToken.NewToken,
	)

	tokenDependencies = fx.Provide(
		func(r repoNonce.Repository) usecaseToken.NonceRepository { return r },
		func(r repoTransaction.Repository) usecaseToken.TransactionRepository { return r },
		func(c *ethereumpkg.Client) usecaseToken.EthereumClient { return c },
		func(t *contracts.CashbackToken) usecaseToken.TokenContract { return t },
	)
)
