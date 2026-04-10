package apperror

import (
	"errors"
	"fmt"
)

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

func As(err error, code string) bool {
	if appErr, ok := errors.AsType[AppError](err); ok {
		return appErr.Code == code
	}
	return false
}
