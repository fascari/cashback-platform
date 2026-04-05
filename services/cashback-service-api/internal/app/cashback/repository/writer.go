package repository

import (
	"context"

	"github.com/cashback-platform/kit/gormtx"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/domain"
)

func (r Repository) Create(ctx context.Context, cashback domain.Cashback) (domain.Cashback, error) {
	model := new(fromDomain(cashback))

	db := r.db.WithContext(ctx)
	if tx := gormtx.ExtractTx(ctx); tx != nil {
		db = tx.WithContext(ctx)
	}

	if err := db.Create(model).Error; err != nil {
		return domain.Cashback{}, err
	}

	return model.toDomain(), nil
}

func (r Repository) Update(ctx context.Context, cashback domain.Cashback) error {
	return r.db.WithContext(ctx).Save(new(fromDomain(cashback))).Error
}
