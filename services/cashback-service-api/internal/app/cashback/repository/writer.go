package repository

import (
	"context"

	"github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/domain"
)

func (r Repository) Create(ctx context.Context, cashback domain.Cashback) (domain.Cashback, error) {
	model := fromDomain(cashback)

	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return domain.Cashback{}, err
	}

	return model.toDomain(), nil
}

func (r Repository) Update(ctx context.Context, cashback domain.Cashback) error {
	model := fromDomain(cashback)
	return r.db.WithContext(ctx).Save(&model).Error
}
