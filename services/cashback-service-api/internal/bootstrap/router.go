package bootstrap

import (
	"net/http"

	"github.com/cashback-platform/services/cashback-service-api/internal/middleware"

	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
)

const serviceName = "cashback-service-api"

var Router = fx.Module("router",
	fx.Provide(NewRouters),
)

type RouterOut struct {
	fx.Out

	MainRouter *chi.Mux   `name:"main"`
	APIRouter  chi.Router `name:"api"`
}

func NewRouters() RouterOut {
	mainRouter := chi.NewRouter()

	mainRouter.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	var apiRouter chi.Router
	mainRouter.Route("/api/v1", func(r chi.Router) {
		middleware.Setup(r, serviceName)
		apiRouter = r
	})

	return RouterOut{
		MainRouter: mainRouter,
		APIRouter:  apiRouter,
	}
}
