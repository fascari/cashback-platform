package modules

import (
	"github.com/cashback-platform/services/cashback-service-api/internal/app/user/handler/createuser"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/user/handler/finduser"
	userrepo "github.com/cashback-platform/services/cashback-service-api/internal/app/user/repository"
	createuseruc "github.com/cashback-platform/services/cashback-service-api/internal/app/user/usecase/createuser"
	finduseruc "github.com/cashback-platform/services/cashback-service-api/internal/app/user/usecase/finduser"

	"go.uber.org/fx"
)

var (
	userFactories = fx.Provide(
		userrepo.New,
		createuseruc.New,
		finduseruc.New,
		createuser.NewHandler,
		finduser.NewHandler,
	)

	userDependencies = fx.Provide(
		func(repo userrepo.Repository) createuseruc.Repository {
			return repo
		},
		func(repo userrepo.Repository) finduseruc.Repository {
			return repo
		},
	)

	userInvokes = fx.Invoke(
		func(params RouterParams, h createuser.Handler) {
			createuser.RegisterEndpoint(params.APIRouter, h)
		},
		func(params RouterParams, h finduser.Handler) {
			finduser.RegisterEndpoint(params.APIRouter, h)
		},
	)

	User = fx.Options(
		userFactories,
		userDependencies,
		userInvokes,
	)
)
