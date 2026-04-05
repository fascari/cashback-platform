package createuser

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	userhandler "github.com/cashback-platform/services/cashback-service-api/internal/app/user/handler"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/user/usecase/createuser"
	"github.com/cashback-platform/services/cashback-service-api/pkg/errorhandler"
	"github.com/cashback-platform/services/cashback-service-api/pkg/httpjson"
)

const Path = "/users"

type Handler struct {
	useCase createuser.UseCase
}

func NewHandler(useCase createuser.UseCase) Handler {
	return Handler{
		useCase: useCase,
	}
}

func RegisterEndpoint(r chi.Router, h Handler) {
	r.Post(Path, h.Handle)
}

func (h Handler) Handle(w http.ResponseWriter, r *http.Request) {
	var payload InputPayload
	if err := httpjson.Read(r, &payload); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	if err := payload.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.useCase.Execute(r.Context(), payload.ExternalID, payload.Email, payload.WalletAddress)
	if err != nil {
		errorhandler.Render(w, err, userhandler.ErrorMapping)
		return
	}

	httpjson.Write(w, http.StatusCreated, ToOutputPayload(user))
}
