package repository

import (
	"context"

	"github.com/cashback-platform/services/cashback-service-api/internal/app/user/domain"
)

func (r Repository) Create(ctx context.Context, user domain.User) (domain.User, error) {
	model := fromDomain(user)

	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return domain.User{}, err
	}

	return model.toDomain(), nil
}
