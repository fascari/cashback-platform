package repository

import (
	"context"
	"errors"

	"github.com/cashback-platform/services/cashback-service-api/internal/app/user/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r Repository) FindByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	var user userModel

	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, errors.New("user not found")
		}
		return domain.User{}, err
	}

	return user.toDomain(), nil
}

func (r Repository) FindByExternalID(ctx context.Context, externalID string) (domain.User, error) {
	var user userModel

	err := r.db.WithContext(ctx).Where("external_id = ?", externalID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, errors.New("user not found")
		}
		return domain.User{}, err
	}

	return user.toDomain(), nil
}

func (r Repository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	var user userModel

	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, errors.New("user not found")
		}
		return domain.User{}, err
	}

	return user.toDomain(), nil
}
