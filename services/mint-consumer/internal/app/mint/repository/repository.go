package repository

import (
	"context"

	"github.com/cashback-platform/kit/gormtx"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) Repository {
	return Repository{db: db}
}

func (r Repository) conn(ctx context.Context) *gorm.DB {
	if tx := gormtx.ExtractTx(ctx); tx != nil {
		return tx.WithContext(ctx)
	}
	return r.db.WithContext(ctx)
}
