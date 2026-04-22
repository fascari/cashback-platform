package modules

import (
	"go.uber.org/fx"

	cashbackrepo "github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/repository"
	deposithandler "github.com/cashback-platform/services/cashback-service-api/internal/app/deposit/handler"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/deposit/handler/depositdetected"
	depositrepo "github.com/cashback-platform/services/cashback-service-api/internal/app/deposit/repository"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/deposit/usecase/processdeposit"
	userrepo "github.com/cashback-platform/services/cashback-service-api/internal/app/user/repository"
)

var (
	depositFactories = fx.Provide(
		depositrepo.New,
		processdeposit.New,
		depositdetected.New,
		deposithandler.NewDepositConsumer,
	)

	depositDependencies = fx.Provide(
		func(repo userrepo.Repository) processdeposit.UserRepository {
			return repo
		},
		func(repo depositrepo.Repository) processdeposit.DepositRepository {
			return repo
		},
		func(repo cashbackrepo.Repository) processdeposit.CashbackRepository {
			return repo
		},
	)

	depositInvokes = fx.Invoke(deposithandler.StartConsumer)

	Deposit = fx.Options(
		depositFactories,
		depositDependencies,
		depositInvokes,
	)
)
