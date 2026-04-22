package depositdetected

import (
	"context"
	"errors"
	"fmt"

	"github.com/cashback-platform/kit/apperror"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/deposit/domain"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/deposit/usecase/processdeposit"
)

type Handler struct {
	useCase processdeposit.UseCase
}

func New(uc processdeposit.UseCase) Handler {
	return Handler{useCase: uc}
}

// Returns a permanentError for payloads that should never be retried (parse errors, duplicate deposits).
// Returns a transient error for infrastructure failures that may resolve on retry.
func (h Handler) Handle(ctx context.Context, data []byte) error {
	input, err := parsePayload(data)
	if err != nil {
		return permanentError{cause: fmt.Errorf("invalid payload: %w", err)}
	}

	if err := h.useCase.Execute(ctx, input); err != nil {
		if apperror.As(err, domain.ErrCodeDepositAlreadyProcessed) {
			return permanentError{cause: err}
		}
		if apperror.As(err, domain.ErrCodeDepositInvalidUser) {
			return permanentError{cause: err}
		}
		return err
	}

	return nil
}

// permanentError signals that the NATS message should be Term'd rather than NAK'd.
type permanentError struct {
	cause error
}

func (e permanentError) Error() string { return e.cause.Error() }
func (e permanentError) Unwrap() error { return e.cause }

// IsPermanent reports whether the error is a permanent failure (Term the message).
func IsPermanent(err error) bool {
	var p permanentError
	return errors.As(err, &p)
}
