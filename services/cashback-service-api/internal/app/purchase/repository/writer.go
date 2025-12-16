package repository

import (
	"context"

	"github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/domain"
)

func (r Repository) Create(ctx context.Context, purchase domain.Purchase) (domain.Purchase, error) {
	model := fromDomain(purchase)

	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return domain.Purchase{}, err
	}

	return model.toDomain(), nil
}
