package handler

import (
	"net/http"

	"github.com/cashback-platform/kit/errorhandler"
	cashdomain "github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/domain"
	calculatecashbackuc "github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/usecase/calculatecashback"
)

var ErrorMapping = errorhandler.ErrorMapping{
	cashdomain.ErrCodeCashbackNotFound:      http.StatusNotFound,
	cashdomain.ErrCodeCashbackAlreadyExists: http.StatusConflict,

	calculatecashbackuc.ErrCodePurchaseNotFound: http.StatusNotFound,
	calculatecashbackuc.ErrCodeUserNotFound:     http.StatusNotFound,
}
