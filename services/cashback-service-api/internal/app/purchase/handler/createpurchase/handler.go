package createpurchase

import (
	"net/http"
	"strconv"

	purchasehandler "github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/handler"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/usecase/createpurchase"
	"github.com/cashback-platform/services/cashback-service-api/pkg/errorhandler"
	httpjson "github.com/cashback-platform/services/cashback-service-api/pkg/http"

	"github.com/go-chi/chi/v5"
)

const Path = "/purchases"

type Handler struct {
	useCase createpurchase.UseCase
}

func NewHandler(useCase createpurchase.UseCase) Handler {
	return Handler{
		useCase: useCase,
	}
}

func RegisterEndpoint(r chi.Router, h Handler) {
	r.Post(Path, h.Handle)
}

func (h Handler) Handle(w http.ResponseWriter, r *http.Request) {
	var payload InputPayload
	if err := httpjson.ReadJSON(r, &payload); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	if err := payload.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(payload.UserID, 10, 64)
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	purchase, err := h.useCase.Execute(r.Context(), userID, payload.Amount, payload.Merchant)
	if err != nil {
		errorhandler.Render(w, err, purchasehandler.ErrorMapping)
		return
	}

	httpjson.WriteJSON(w, http.StatusCreated, ToOutputPayload(purchase))
}
