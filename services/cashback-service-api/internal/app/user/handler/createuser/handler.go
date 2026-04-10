package createuser

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/cashback-platform/kit/errorhandler"
	"github.com/cashback-platform/kit/httpjson"
	userhandler "github.com/cashback-platform/services/cashback-service-api/internal/app/user/handler"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/user/usecase/createuser"
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
	payload := new(InputPayload)
	if err := httpjson.Read(r, payload); err != nil {
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
