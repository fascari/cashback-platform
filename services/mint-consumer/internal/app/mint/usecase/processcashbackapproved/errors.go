package processcashbackapproved

import "github.com/cashback-platform/kit/apperror"

const ErrCodeInvalidPayload = "error_process_cashback_approved_invalid_payload"

// ErrInvalidPayload is returned when the NATS message cannot be JSON-decoded.
var ErrInvalidPayload = apperror.New(ErrCodeInvalidPayload, "invalid cashback approved payload")
