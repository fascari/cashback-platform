// Package errorhandler provides utilities for handling HTTP errors consistently.
package errorhandler

import (
	"errors"
	"net/http"
)

// Common HTTP errors
var (
	ErrBadRequest          = NewHTTPError(http.StatusBadRequest, "bad request")
	ErrUnauthorized        = NewHTTPError(http.StatusUnauthorized, "unauthorized")
	ErrForbidden           = NewHTTPError(http.StatusForbidden, "forbidden")
	ErrNotFound            = NewHTTPError(http.StatusNotFound, "not found")
	ErrConflict            = NewHTTPError(http.StatusConflict, "conflict")
	ErrUnprocessableEntity = NewHTTPError(http.StatusUnprocessableEntity, "unprocessable entity")
	ErrInternalServer      = NewHTTPError(http.StatusInternalServerError, "internal server error")
)

// HTTPError represents an error that has an associated HTTP status code.
type HTTPError struct {
	Code    int
	Message string
	Err     error
}

func (e *HTTPError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *HTTPError) Unwrap() error {
	return e.Err
}

// NewHTTPError creates a new HTTPError with the given status code and message.
func NewHTTPError(code int, message string) *HTTPError {
	return &HTTPError{
		Code:    code,
		Message: message,
	}
}

// WrapHTTPError wraps an existing error with an HTTP status code and message.
func WrapHTTPError(code int, message string, err error) *HTTPError {
	return &HTTPError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Render writes the HTTP error response to the ResponseWriter.
// It determines the status code and message based on the error type.
func Render(w http.ResponseWriter, err error) {
	var httpErr *HTTPError
	if errors.As(err, &httpErr) {
		http.Error(w, httpErr.Message, httpErr.Code)
		return
	}

	// Default to internal server error for unknown errors
	http.Error(w, "internal server error", http.StatusInternalServerError)
}

// RenderWithCode writes the HTTP error response with a specific status code.
func RenderWithCode(w http.ResponseWriter, code int, message string) {
	http.Error(w, message, code)
}

// CodeFromError extracts the HTTP status code from an error.
// Returns 500 if the error is not an HTTPError.
func CodeFromError(err error) int {
	var httpErr *HTTPError
	if errors.As(err, &httpErr) {
		return httpErr.Code
	}
	return http.StatusInternalServerError
}
