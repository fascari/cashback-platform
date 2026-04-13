package repository

import (
	"context"
	"encoding/json"

	"github.com/cashback-platform/kit/events"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/cashback/domain"
	"gorm.io/gorm"
)

func (r Repository) Create(ctx context.Context, cashback domain.Cashback) (domain.Cashback, error) {
	model := new(fromDomain(cashback))
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return domain.Cashback{}, err
	}
	return model.toDomain(), nil
}

func (r Repository) CreateWithEvent(ctx context.Context, cashback domain.Cashback, buildPayload func(domain.Cashback) any) (domain.Cashback, error) {
	m := fromDomain(cashback)
	var created cashbackModel
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&m).Error; err != nil {
			return err
		}
		created = m

		payload, err := json.Marshal(buildPayload(m.toDomain()))
		if err != nil {
			return err
		}
		return r.outboxWriter.CreateWithTx(tx, events.CashbackApproved, "cashback", m.ID, payload)
	})
	if err != nil {
		return domain.Cashback{}, err
	}
	return created.toDomain(), nil
}

func (r Repository) Update(ctx context.Context, cashback domain.Cashback) error {
	return r.db.WithContext(ctx).Save(new(fromDomain(cashback))).Error
}
