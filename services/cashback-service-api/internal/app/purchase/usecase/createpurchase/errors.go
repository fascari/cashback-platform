package createpurchase

import "errors"

var (
	ErrInvalidAmount   = errors.New("invalid purchase amount")
	ErrInvalidUserID   = errors.New("invalid user ID")
	ErrInvalidMerchant = errors.New("invalid merchant ID")
)
