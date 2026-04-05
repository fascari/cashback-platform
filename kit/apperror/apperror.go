package apperror

import (
	"errors"
	"fmt"
)

// AppError is a typed error carrying a machine-readable code and a human-readable message.
// It travels through all application layers and is used by errorhandler to resolve HTTP status codes.
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func New(code string, format string, args ...any) AppError {
	return AppError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

func (e AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// As reports whether any error in err's chain is an AppError with the given code.
func As(err error, code string) bool {
	var appErr AppError
	if errors.As(err, &appErr) {
		return appErr.Code == code
	}
	return false
}
