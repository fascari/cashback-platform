package modules

import (
	"github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/handler/calculatecashback"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/handler/findusercashback"
	cashbackrepo "github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/repository"
	calculatecashbackuc "github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/usecase/calculatecashback"
	findusercashbackuc "github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/usecase/findusercashback"
	purchaserepo "github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/repository"
	userrepo "github.com/cashback-platform/services/cashback-service-api/internal/app/user/repository"
	outboxrepo "github.com/cashback-platform/services/cashback-service-api/internal/infra/messaging/outbox/repository"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

var (
	cashbackFactories = fx.Provide(
		func(db *gorm.DB, outboxRepo outboxrepo.Repository) cashbackrepo.Repository {
			return cashbackrepo.New(db, outboxRepo)
		},
		calculatecashbackuc.New,
		findusercashbackuc.New,
		calculatecashback.NewHandler,
		findusercashback.NewHandler,
	)

	cashbackDependencies = fx.Provide(
		func(repo cashbackrepo.Repository) calculatecashbackuc.Repository {
			return repo
		},
		func(repo purchaserepo.Repository) calculatecashbackuc.PurchaseRepository {
			return repo
		},
		func(repo userrepo.Repository) calculatecashbackuc.UserRepository {
			return repo
		},
		func(repo cashbackrepo.Repository) findusercashbackuc.Repository {
			return repo
		},
	)

	cashbackInvokes = fx.Invoke(
		func(params RouterParams, h calculatecashback.Handler) {
			calculatecashback.RegisterEndpoint(params.APIRouter, h)
		},
		func(params RouterParams, h findusercashback.Handler) {
			findusercashback.RegisterEndpoint(params.APIRouter, h)
		},
	)

	Cashback = fx.Options(
		cashbackFactories,
		cashbackDependencies,
		cashbackInvokes,
	)
)
