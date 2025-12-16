package repository

import (
	"context"
	"errors"

	"github.com/cashback-platform/services/cashback-service-api/internal/app/purchase/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r Repository) FindByID(ctx context.Context, id uuid.UUID) (domain.Purchase, error) {
	var purchase purchaseModel

	err := r.db.WithContext(ctx).First(&purchase, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Purchase{}, errors.New("purchase not found")
		}
		return domain.Purchase{}, err
	}

	return purchase.toDomain(), nil
}

func (r Repository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Purchase, error) {
	var purchases []purchaseModel

	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&purchases).Error
	if err != nil {
		return nil, err
	}

	result := make([]domain.Purchase, len(purchases))
	for i, p := range purchases {
		result[i] = p.toDomain()
	}

	return result, nil
}
