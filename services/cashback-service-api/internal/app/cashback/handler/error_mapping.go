package handler

import (
	"net/http"

	cashdomain "github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/domain"
	calculatecashbackuc "github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/usecase/calculatecashback"
	"github.com/cashback-platform/services/cashback-service-api/pkg/errorhandler"
)

// ErrorMapping maps cashback domain and usecase error codes to HTTP status codes.
var ErrorMapping = errorhandler.ErrorMapping{
	cashdomain.ErrCodeCashbackNotFound:      http.StatusNotFound,
	cashdomain.ErrCodeCashbackAlreadyExists: http.StatusConflict,

	calculatecashbackuc.ErrCodePurchaseNotFound: http.StatusNotFound,
	calculatecashbackuc.ErrCodeUserNotFound:     http.StatusNotFound,
}
