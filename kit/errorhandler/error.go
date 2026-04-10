package errorhandler

import (
	"errors"
	"net/http"

	"github.com/cashback-platform/kit/apperror"
	"github.com/cashback-platform/kit/httpjson"
	"github.com/cashback-platform/kit/logger"
)

type ErrorMapping map[string]int

func Render(w http.ResponseWriter, err error, mapping ...ErrorMapping) {
	var m ErrorMapping
	if len(mapping) > 0 {
		m = mapping[0]
	}
	httpCode := resolveHTTPCode(err, m)
	logError(httpCode, err)
	httpjson.Error(w, httpCode, err)
}

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
