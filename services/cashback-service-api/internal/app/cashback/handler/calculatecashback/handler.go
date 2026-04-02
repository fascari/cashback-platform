package calculatecashback

import (
	"net/http"
	"strconv"

	cashbackhandler "github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/handler"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/usecase/calculatecashback"
	"github.com/cashback-platform/services/cashback-service-api/pkg/apperror"
	"github.com/cashback-platform/services/cashback-service-api/pkg/errorhandler"
	httpjson "github.com/cashback-platform/services/cashback-service-api/pkg/http"
	"github.com/go-chi/chi/v5"
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	purchaseID, err := strconv.ParseInt(payload.PurchaseID, 10, 64)
	if err != nil {
		errorhandler.RenderWithCode(w, http.StatusBadRequest, "invalid purchase ID")
		return
	}

	cashback, err := h.useCase.Execute(r.Context(), purchaseID)
	if err != nil {
		if apperror.As(err, calculatecashback.ErrCodeFailedToPublishEvent) {
			httpjson.WriteJSON(w, http.StatusCreated, ToOutputPayload(cashback))
			return
		}

		errorhandler.Render(w, err, cashbackhandler.ErrorMapping)
		return
	}

	httpjson.WriteJSON(w, http.StatusCreated, ToOutputPayload(cashback))
}
