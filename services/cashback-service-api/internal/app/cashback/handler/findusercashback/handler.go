package findusercashback

import (
	"net/http"

	"github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/usecase/findusercashback"
	httpjson "github.com/cashback-platform/services/cashback-service-api/pkg/http"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

const Path = "/users/{user_id}/cashback"

type Handler struct {
	useCase findusercashback.UseCase
}

func NewHandler(useCase findusercashback.UseCase) Handler {
	return Handler{
		useCase: useCase,
	}
}

func RegisterEndpoint(r chi.Router, h Handler) {
	r.Get(Path, h.Handle)
}

func (h Handler) Handle(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	summary, err := h.useCase.Execute(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, ToOutputPayload(summary))
}
