package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/cashback-platform/services/cashback-service-api/pkg/apperror"
	"github.com/cashback-platform/services/cashback-service-api/pkg/logger"
)

func WriteJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		logger.Error("failed to encode JSON response", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func ReadJSON(r *http.Request, payload any) error {
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		logger.Error("failed to decode JSON request", "error", err)
		return err
	}
	return nil
}

func WriteError(w http.ResponseWriter, statusCode int, err error) {
	statusText := strings.ToUpper(strings.ReplaceAll(http.StatusText(statusCode), " ", "_"))
	appErr := apperror.New(statusText, "%s", err.Error())
	WriteJSON(w, statusCode, appErr)
}
