package handler

import (
	"net/http"

	userdomain "github.com/cashback-platform/services/cashback-service-api/internal/app/user/domain"
	"github.com/cashback-platform/kit/errorhandler"
)

var ErrorMapping = errorhandler.ErrorMapping{
	userdomain.ErrCodeUserNotFound:      http.StatusNotFound,
	userdomain.ErrCodeUserAlreadyExists: http.StatusConflict,
}
