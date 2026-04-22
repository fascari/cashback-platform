package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/cashback-platform/services/cashback-service-api/internal/app/deposit/domain"
)

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) Repository {
	return Repository{db: db}
}

func (r Repository) Save(ctx context.Context, receipt domain.DepositReceipt) (domain.DepositReceipt, error) {
	m := fromDomain(receipt)
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return domain.DepositReceipt{}, err
	}
	return m.toDomain(), nil
}

func (r Repository) ExistsByTxHash(ctx context.Context, txHash string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&depositReceiptModel{}).
		Where("tx_hash = ?", txHash).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
