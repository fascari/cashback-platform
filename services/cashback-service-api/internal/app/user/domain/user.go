package domain

import (
	"time"

	"github.com/cashback-platform/kit/apperror"
)

const (
	ErrCodeUserNotFound      = "error_user_not_found"
	ErrCodeUserAlreadyExists = "error_user_already_exists"
)

var (
	ErrUserNotFound      = apperror.New(ErrCodeUserNotFound, "user not found")
	ErrUserAlreadyExists = apperror.New(ErrCodeUserAlreadyExists, "user already exists")
)

type User struct {
	ID            int64
	ExternalID    string
	Email         string
	WalletAddress string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
