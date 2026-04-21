package handler

import (
	"net/http"

	"github.com/cashback-platform/kit/errorhandler"
	userdomain "github.com/cashback-platform/services/cashback-service-api/internal/app/user/domain"
)

var ErrorMapping = errorhandler.ErrorMapping{
	userdomain.ErrCodeUserNotFound:      http.StatusNotFound,
	userdomain.ErrCodeUserAlreadyExists: http.StatusConflict,
}
