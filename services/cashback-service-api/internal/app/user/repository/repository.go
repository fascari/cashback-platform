// Package repository implements data persistence for user entities.
package repository

import (
	"gorm.io/gorm"
)

// Repository handles user data persistence operations.
// It provides methods for both reading and writing user records.
type Repository struct {
	db *gorm.DB
}

// New creates a new user repository instance.
func New(db *gorm.DB) Repository {
	return Repository{
		db: db,
	}
}
