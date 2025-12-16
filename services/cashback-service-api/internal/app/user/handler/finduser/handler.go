package finduser

import (
	"net/http"

	"github.com/cashback-platform/services/cashback-service-api/internal/app/user/usecase/finduser"
	httpjson "github.com/cashback-platform/services/cashback-service-api/pkg/http"
	"github.com/google/uuid"

	"github.com/go-chi/chi/v5"
)

const Path = "/users/{id}"

type Handler struct {
	useCase finduser.UseCase
}

func NewHandler(useCase finduser.UseCase) Handler {
	return Handler{
		useCase: useCase,
	}
}

func RegisterEndpoint(r chi.Router, h Handler) {
	r.Get(Path, h.Handle)
}

func (h Handler) Handle(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	user, err := h.useCase.Execute(r.Context(), id)
	if err != nil {
		if err.Error() == "user not found" {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, ToOutputPayload(user))
}
