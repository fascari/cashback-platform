package modules

import (
	repoNonce "github.com/cashback-platform/services/blockchain-adapter/internal/repository/nonce"
	repoTransaction "github.com/cashback-platform/services/blockchain-adapter/internal/repository/transaction"
	usecaseToken "github.com/cashback-platform/services/blockchain-adapter/internal/usecase"
	"go.uber.org/fx"
)

var Token = fx.Options(
	fx.Provide(
		repoTransaction.NewRepository,
		repoNonce.NewRepository,
		usecaseToken.NewToken,
	),
)
