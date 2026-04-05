package repository

import (
	"context"

	"github.com/cashback-platform/kit/gormtx"
	"gorm.io/gorm"
)

// Repository provides data access for mint requests and processed events.
type Repository struct {
	db *gorm.DB
}

// New creates a Repository backed by the given database connection.
func New(db *gorm.DB) Repository {
	return Repository{db: db}
}

// conn returns the active transaction from ctx if one exists, otherwise the plain connection.
func (r Repository) conn(ctx context.Context) *gorm.DB {
	if tx := gormtx.ExtractTx(ctx); tx != nil {
		return tx.WithContext(ctx)
	}
	return r.db.WithContext(ctx)
}
