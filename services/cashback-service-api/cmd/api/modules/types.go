package modules

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
)

type RouterParams struct {
	fx.In

	Router    *chi.Mux   `name:"main"`
	APIRouter chi.Router `name:"api"`
}
