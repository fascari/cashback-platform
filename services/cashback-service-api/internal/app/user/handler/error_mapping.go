package handler

import (
	"net/http"

	userdomain "github.com/cashback-platform/services/cashback-service-api/internal/app/user/domain"
	"github.com/cashback-platform/services/cashback-service-api/pkg/errorhandler"
)

// ErrorMapping maps user domain error codes to HTTP status codes.
var ErrorMapping = errorhandler.ErrorMapping{
	userdomain.ErrCodeUserNotFound:      http.StatusNotFound,
	userdomain.ErrCodeUserAlreadyExists: http.StatusConflict,
}
