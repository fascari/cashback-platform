package calculatecashback

import "github.com/cashback-platform/kit/apperror"

const (
	ErrCodePurchaseNotFound = "error_calculatecashback_purchase_not_found"
	ErrCodeUserNotFound     = "error_calculatecashback_user_not_found"
)

var (
	ErrPurchaseNotFound = apperror.New(ErrCodePurchaseNotFound, "purchase not found")
	ErrUserNotFound     = apperror.New(ErrCodeUserNotFound, "user not found")
)
