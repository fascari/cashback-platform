package balance

import (
	"net/http"
	"strconv"

	"github.com/cashback-platform/kit/errorhandler"
	"github.com/cashback-platform/kit/httpjson"
	userhandler "github.com/cashback-platform/services/cashback-service-api/internal/app/user/handler"
	balanceuc "github.com/cashback-platform/services/cashback-service-api/internal/app/user/usecase/balance"

	"github.com/go-chi/chi/v5"
)

const Path = "/users/{id}/balance"

type Handler struct {
	useCase balanceuc.UseCase
}

func NewHandler(useCase balanceuc.UseCase) Handler {
	return Handler{useCase: useCase}
}

func RegisterEndpoint(r chi.Router, h Handler) {
	r.Get(Path, h.Handle)
}

func (h Handler) Handle(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	output, err := h.useCase.Execute(r.Context(), id)
	if err != nil {
		errorhandler.Render(w, err, userhandler.ErrorMapping)
		return
	}

	httpjson.Write(w, http.StatusOK, ToOutputPayload(output))
}
