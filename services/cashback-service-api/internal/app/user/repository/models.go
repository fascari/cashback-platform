package repository

import (
	"time"

	"github.com/cashback-platform/services/cashback-service-api/internal/app/user/domain"
)

type userModel struct {
	ID            int64  `gorm:"primaryKey;autoIncrement"`
	ExternalID    string `gorm:"uniqueIndex;not null"`
	Email         string `gorm:"uniqueIndex;not null"`
	WalletAddress string `gorm:"not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (userModel) TableName() string {
	return "cashback.users"
}

func (m userModel) toDomain() domain.User {
	return domain.User{
		ID:            m.ID,
		ExternalID:    m.ExternalID,
		Email:         m.Email,
		WalletAddress: m.WalletAddress,
		CreatedAt:     m.CreatedAt.UTC(),
		UpdatedAt:     m.UpdatedAt.UTC(),
	}
}

func fromDomain(u domain.User) userModel {
	return userModel{
		ID:            u.ID,
		ExternalID:    u.ExternalID,
		Email:         u.Email,
		WalletAddress: u.WalletAddress,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}
}
