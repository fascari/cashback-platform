package calculatecashback

import (
	"net/http"

	"github.com/cashback-platform/services/cashback-service-api/pkg/errorhandler"
)

var (
	ErrPurchaseNotFound      = errorhandler.NewHTTPError(http.StatusNotFound, "purchase not found")
	ErrUserNotFound          = errorhandler.NewHTTPError(http.StatusNotFound, "user not found")
	ErrCashbackAlreadyExists = errorhandler.NewHTTPError(http.StatusConflict, "cashback already exists for this purchase")
	ErrFailedToPublishEvent  = errorhandler.NewHTTPError(http.StatusCreated, "cashback created but event publishing failed")
	ErrInvalidPurchaseID     = errorhandler.NewHTTPError(http.StatusBadRequest, "invalid purchase ID")
)
