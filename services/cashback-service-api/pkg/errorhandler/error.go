package errorhandler

import (
	"errors"
	"net/http"

	"github.com/cashback-platform/kit/apperror"
	"github.com/cashback-platform/kit/logger"
	"github.com/cashback-platform/services/cashback-service-api/pkg/httpjson"
)

// ErrorMapping maps apperror codes to HTTP status codes.
type ErrorMapping map[string]int

// Render resolves the HTTP status code from err via the provided ErrorMapping,
// logs the error at the appropriate severity, and writes a JSON error response.
// mapping is optional — omit it to fall back to 500 for all errors.
func Render(w http.ResponseWriter, err error, mapping ...ErrorMapping) {
	var m ErrorMapping
	if len(mapping) > 0 {
		m = mapping[0]
	}
	httpCode := resolveHTTPCode(err, m)
	logError(httpCode, err)
	httpjson.Error(w, httpCode, err)
}

// RenderWithCode writes a plain-text error response with an explicit status code.
func RenderWithCode(w http.ResponseWriter, code int, message string) {
	http.Error(w, message, code)
}

func resolveHTTPCode(err error, mapping ErrorMapping) int {
	if appErr, ok := errors.AsType[apperror.AppError](err); ok && len(mapping) > 0 {
		if code, ok := mapping[appErr.Code]; ok {
			return code
		}
	}
	return http.StatusInternalServerError
}

func logError(httpCode int, err error) {
	if httpCode >= http.StatusInternalServerError {
		logger.Error("request error", "error", err)
		return
	}
	logger.Warn("request error", "error", err)
}
