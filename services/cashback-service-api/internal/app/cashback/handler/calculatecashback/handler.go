package calculatecashback

import (
	"errors"
	"net/http"

	"github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/usecase/calculatecashback"
	"github.com/cashback-platform/services/cashback-service-api/pkg/errorhandler"
	httpjson "github.com/cashback-platform/services/cashback-service-api/pkg/http"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

const Path = "/cashback/calculate"

type Handler struct {
	useCase calculatecashback.UseCase
}

func NewHandler(useCase calculatecashback.UseCase) Handler {
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
		errorhandler.RenderWithCode(w, http.StatusBadRequest, "invalid payload")
		return
	}

	if err := payload.Validate(); err != nil {
		errorhandler.Render(w, err)
		return
	}

	purchaseID, err := uuid.Parse(payload.PurchaseID)
	if err != nil {
		errorhandler.Render(w, calculatecashback.ErrInvalidPurchaseID)
		return
	}

	cashback, err := h.useCase.Execute(r.Context(), purchaseID)
	if err != nil {
		if errors.Is(err, calculatecashback.ErrFailedToPublishEvent) {
			httpjson.WriteJSON(w, http.StatusCreated, ToOutputPayload(cashback))
			return
		}

		errorhandler.Render(w, err)
		return
	}

	httpjson.WriteJSON(w, http.StatusCreated, ToOutputPayload(cashback))
}
