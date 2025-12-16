// Package domain contains the core business entities and rules for users.
package domain

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system.
// Each user has a unique external ID, email, and blockchain wallet address.
type User struct {
	ID            uuid.UUID
	ExternalID    string
	Email         string
	WalletAddress string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
