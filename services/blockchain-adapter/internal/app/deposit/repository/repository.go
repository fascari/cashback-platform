package repository

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/cashback-platform/services/blockchain-adapter/internal/app/deposit/domain"
)

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) Repository {
	return Repository{db: db}
}

func (r Repository) Save(ctx context.Context, d domain.Deposit) error {
	m := fromDomain(d)
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&m).Error
}

func (r Repository) MaxBlockNumber(ctx context.Context, chainID string) (int64, error) {
	var result int64
	err := r.db.WithContext(ctx).
		Model(&depositModel{}).
		Where("chain_id = ?", chainID).
		Select("COALESCE(MAX(block_number), 0)").
		Scan(&result).Error
	return result, err
}
