package validator

import "errors"

var (
	ErrInvalidEmail         = errors.New("invalid email")
	ErrInvalidWalletAddress = errors.New("invalid wallet address")
	ErrRequired             = errors.New("required field")
)

func ValidateEmail(email string) error {
	if email == "" {
		return ErrRequired
	}
	if len(email) < 3 || !containsAtSymbol(email) {
		return ErrInvalidEmail
	}
	return nil
}

func ValidateWalletAddress(address string) error {
	if address == "" {
		return ErrRequired
	}
	if len(address) < 20 {
		return ErrInvalidWalletAddress
	}
	return nil
}

func containsAtSymbol(s string) bool {
	for _, c := range s {
		if c == '@' {
			return true
		}
	}
	return false
}
