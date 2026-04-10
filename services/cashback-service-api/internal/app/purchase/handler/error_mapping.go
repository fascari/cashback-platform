package handler

import (
	"net/http"

	"github.com/cashback-platform/kit/errorhandler"
	purchasedomain "github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/domain"
	createpurchaseuc "github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/usecase/createpurchase"
)

// ErrorMapping maps purchase domain and usecase error codes to HTTP status codes.
var ErrorMapping = errorhandler.ErrorMapping{
	purchasedomain.ErrCodePurchaseNotFound: http.StatusNotFound,

	createpurchaseuc.ErrCodeInvalidAmount:   http.StatusUnprocessableEntity,
	createpurchaseuc.ErrCodeInvalidUserID:   http.StatusUnprocessableEntity,
	createpurchaseuc.ErrCodeInvalidMerchant: http.StatusUnprocessableEntity,
}
