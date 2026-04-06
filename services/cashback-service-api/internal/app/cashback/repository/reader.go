package repository

import (
	"context"
	"errors"

	"github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/domain"
	"gorm.io/gorm"
)

func (r Repository) FindByID(ctx context.Context, id int64) (domain.Cashback, error) {
	cashback := new(cashbackModel)

	err := r.db.WithContext(ctx).First(cashback, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Cashback{}, domain.ErrCashbackNotFound
		}
		return domain.Cashback{}, err
	}

	return cashback.toDomain(), nil
}

func (r Repository) FindByUserID(ctx context.Context, userID int64) ([]domain.Cashback, error) {
	var cashbacks []cashbackModel

	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&cashbacks).Error
	if err != nil {
		return nil, err
	}

	result := make([]domain.Cashback, len(cashbacks))
	for i, c := range cashbacks {
		result[i] = c.toDomain()
	}

	return result, nil
}

func (r Repository) FindByPurchaseID(ctx context.Context, purchaseID int64) (domain.Cashback, error) {
	cashback := new(cashbackModel)

	err := r.db.WithContext(ctx).
		Where("purchase_id = ?", purchaseID).
		First(cashback).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Cashback{}, domain.ErrCashbackNotFound
		}
		return domain.Cashback{}, err
	}

	return cashback.toDomain(), nil
}

func (r Repository) TotalByUserID(ctx context.Context, userID int64) (float64, error) {
	total := new(float64)
	err := r.db.WithContext(ctx).
		Model(&cashbackModel{}).
		Where("user_id = ? AND status = ?", userID, domain.StatusMinted).
		Select("COALESCE(SUM(amount), 0)").
		Scan(total).Error

	return *total, err
}
