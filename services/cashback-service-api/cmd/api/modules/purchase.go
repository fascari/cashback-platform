package modules

import (
	"github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/handler/createpurchase"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/handler/findpurchase"
	purchaserepo "github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/repository"
	createpurchaseuc "github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/usecase/createpurchase"
	findpurchaseuc "github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/usecase/findpurchase"

	"go.uber.org/fx"
)

var (
	purchaseFactories = fx.Provide(
		purchaserepo.New,
		createpurchaseuc.New,
		findpurchaseuc.New,
		createpurchase.NewHandler,
		findpurchase.NewHandler,
	)

	purchaseDependencies = fx.Provide(
		func(repo purchaserepo.Repository) createpurchaseuc.Repository {
			return repo
		},
		func(repo purchaserepo.Repository) findpurchaseuc.Repository {
			return repo
		},
	)

	purchaseInvokes = fx.Invoke(
		func(params RouterParams, h createpurchase.Handler) {
			createpurchase.RegisterEndpoint(params.APIRouter, h)
		},
		func(params RouterParams, h findpurchase.Handler) {
			findpurchase.RegisterEndpoint(params.APIRouter, h)
		},
	)

	Purchase = fx.Options(
		purchaseFactories,
		purchaseDependencies,
		purchaseInvokes,
	)
)
