package gormtx

import (
	"context"

	"gorm.io/gorm"
)

type (
	ctxKey struct{}

	TransactionManager struct {
		db *gorm.DB
	}
)

func InjectTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, ctxKey{}, tx)
}

func ExtractTx(ctx context.Context) *gorm.DB {
	tx, _ := ctx.Value(ctxKey{}).(*gorm.DB)
	return tx
}

func NewTransactionManager(db *gorm.DB) TransactionManager {
	return TransactionManager{db: db}
}

func (m TransactionManager) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(InjectTx(ctx, tx))
	})
}
