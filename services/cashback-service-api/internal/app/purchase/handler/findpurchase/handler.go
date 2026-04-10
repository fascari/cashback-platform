package findpurchase

import (
	"net/http"
	"strconv"

	purchasehandler "github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/handler"
	findpurchaseuc "github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/usecase/findpurchase"
	"github.com/cashback-platform/kit/errorhandler"
	"github.com/cashback-platform/kit/httpjson"

	"github.com/go-chi/chi/v5"
)

const Path = "/purchases/{id}"

type Handler struct {
	useCase findpurchaseuc.UseCase
}

func NewHandler(useCase findpurchaseuc.UseCase) Handler {
	return Handler{useCase: useCase}
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
		errorhandler.Render(w, err, purchasehandler.ErrorMapping)
		return
	}

	httpjson.Write(w, http.StatusOK, ToOutputPayload(purchase))
}
