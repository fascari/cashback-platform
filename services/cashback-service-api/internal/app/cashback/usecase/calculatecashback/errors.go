package calculatecashback

import "github.com/cashback-platform/services/cashback-service-api/pkg/apperror"

const (
	ErrCodePurchaseNotFound     = "error_calculatecashback_purchase_not_found"
	ErrCodeUserNotFound         = "error_calculatecashback_user_not_found"
	ErrCodeFailedToPublishEvent = "error_calculatecashback_failed_to_publish_event"
)

var (
	ErrPurchaseNotFound     = apperror.New(ErrCodePurchaseNotFound, "purchase not found")
	ErrUserNotFound         = apperror.New(ErrCodeUserNotFound, "user not found")
	ErrFailedToPublishEvent = apperror.New(ErrCodeFailedToPublishEvent, "failed to publish cashback event")
)
