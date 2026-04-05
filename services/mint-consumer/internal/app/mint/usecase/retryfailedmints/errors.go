package retryfailedmints

import "github.com/cashback-platform/kit/apperror"

const ErrCodeFetchFailed = "error_retry_failed_mints_fetch_failed"

// ErrFetchFailed is returned when the repository query for retryable requests fails.
var ErrFetchFailed = apperror.New(ErrCodeFetchFailed, "failed to fetch retryable mint requests")
