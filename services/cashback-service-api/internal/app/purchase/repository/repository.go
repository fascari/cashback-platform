// Package repository implements data persistence for purchase entities.
package repository

import (
	"gorm.io/gorm"
)

// Repository handles purchase data persistence operations.
// It provides methods for both reading and writing purchase records.
type Repository struct {
	db *gorm.DB
}

// New creates a new purchase repository instance.
func New(db *gorm.DB) Repository {
	return Repository{
		db: db,
	}
}
