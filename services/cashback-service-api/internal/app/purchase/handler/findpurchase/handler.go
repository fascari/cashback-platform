package findpurchase

import (
	"net/http"
	"strconv"

	findpurchaseuc "github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/usecase/findpurchase"
	httpjson "github.com/cashback-platform/services/cashback-service-api/pkg/http"

	"github.com/go-chi/chi/v5"
)

const Path = "/purchases/{id}"

type Handler struct {
	useCase findpurchaseuc.UseCase
}

func NewHandler(useCase findpurchaseuc.UseCase) Handler {
	return Handler{
		useCase: useCase,
	}
}

func RegisterEndpoint(r chi.Router, h Handler) {
	r.Get(Path, h.Handle)
}

func (h Handler) Handle(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid purchase id", http.StatusBadRequest)
		return
	}

	purchase, err := h.useCase.Execute(r.Context(), id)
	if err != nil {
		if err.Error() == "purchase not found" {
			http.Error(w, "purchase not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, ToOutputPayload(purchase))
}
