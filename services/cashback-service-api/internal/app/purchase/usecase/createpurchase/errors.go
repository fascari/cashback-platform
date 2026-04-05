package createpurchase

import "github.com/cashback-platform/kit/apperror"

const (
	ErrCodeInvalidAmount   = "error_createpurchase_invalid_amount"
	ErrCodeInvalidUserID   = "error_createpurchase_invalid_user_id"
	ErrCodeInvalidMerchant = "error_createpurchase_invalid_merchant"
)

var (
	ErrInvalidAmount   = apperror.New(ErrCodeInvalidAmount, "invalid purchase amount")
	ErrInvalidUserID   = apperror.New(ErrCodeInvalidUserID, "invalid user ID")
	ErrInvalidMerchant = apperror.New(ErrCodeInvalidMerchant, "invalid merchant ID")
)
