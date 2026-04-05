package gormtx

import (
	"context"

	"gorm.io/gorm"
)

type ctxKey struct{}

// InjectTx stores a GORM transaction in ctx so that repositories can join it.
func InjectTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, ctxKey{}, tx)
}

// ExtractTx retrieves the GORM transaction from ctx, or returns nil if none was set.
func ExtractTx(ctx context.Context) *gorm.DB {
	tx, _ := ctx.Value(ctxKey{}).(*gorm.DB)
	return tx
}

// TransactionManager executes fn inside a GORM database transaction.
// Any tx injected via InjectTx is available to repositories via ExtractTx.
type TransactionManager struct {
	db *gorm.DB
}

// NewTransactionManager returns a TransactionManager backed by db.
func NewTransactionManager(db *gorm.DB) TransactionManager {
	return TransactionManager{db: db}
}

// WithTransaction runs fn inside a database transaction.
// If fn returns an error the transaction is rolled back; otherwise it is committed.
func (m TransactionManager) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(InjectTx(ctx, tx))
	})
}
