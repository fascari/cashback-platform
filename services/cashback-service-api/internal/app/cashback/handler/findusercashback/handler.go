package findusercashback

import (
	"net/http"
	"strconv"

	cashbackhandler "github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/handler"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/usecase/findusercashback"
	"github.com/cashback-platform/kit/errorhandler"
	"github.com/cashback-platform/kit/httpjson"
	"github.com/go-chi/chi/v5"
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
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	summary, err := h.useCase.Execute(r.Context(), userID)
	if err != nil {
		errorhandler.Render(w, err, cashbackhandler.ErrorMapping)
		return
	}

	httpjson.Write(w, http.StatusOK, ToOutputPayload(summary))
}
