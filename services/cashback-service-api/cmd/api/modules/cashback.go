package modules

import (
	"context"

	"github.com/cashback-platform/kit/gormtx"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/handler/calculatecashback"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/handler/findusercashback"
	cashbackrepo "github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/repository"
	calculatecashbackuc "github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/usecase/calculatecashback"
	findusercashbackuc "github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/usecase/findusercashback"
	purchaserepo "github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/repository"
	userrepo "github.com/cashback-platform/services/cashback-service-api/internal/app/user/repository"
	"github.com/cashback-platform/services/cashback-service-api/internal/infra/messaging/outbox"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

var (
	cashbackFactories = fx.Provide(
		cashbackrepo.New,
		calculatecashbackuc.New,
		findusercashbackuc.New,
		calculatecashback.NewHandler,
		findusercashback.NewHandler,
		func(db *gorm.DB) calculatecashbackuc.TransactionManager {
			return gormtx.NewTransactionManager(db)
		},
	)

	cashbackDependencies = fx.Provide(
		func(repo cashbackrepo.Repository) calculatecashbackuc.Repository {
			return repo
		},
		func(repo purchaserepo.Repository) calculatecashbackuc.PurchaseRepository {
			return purchaseRepoAdapter{repo: repo}
		},
		func(repo userrepo.Repository) calculatecashbackuc.UserRepository {
			return userRepoAdapter{repo: repo}
		},
		func(pub outbox.Publisher) calculatecashbackuc.EventPublisher {
			return pub
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

type (
	purchaseRepoAdapter struct {
		repo purchaserepo.Repository
	}

	userRepoAdapter struct {
		repo userrepo.Repository
	}
)

func (a purchaseRepoAdapter) FindByID(ctx context.Context, id int64) (calculatecashbackuc.Purchase, error) {
	p, err := a.repo.FindByID(ctx, id)
	if err != nil {
		return calculatecashbackuc.Purchase{}, err
	}
	return calculatecashbackuc.Purchase{ID: p.ID, UserID: p.UserID, Amount: p.Amount}, nil
}

func (a userRepoAdapter) FindByID(ctx context.Context, id int64) (calculatecashbackuc.User, error) {
	u, err := a.repo.FindByID(ctx, id)
	if err != nil {
		return calculatecashbackuc.User{}, err
	}
	return calculatecashbackuc.User{WalletAddress: u.WalletAddress}, nil
}
