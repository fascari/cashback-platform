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

func (userModel) TableName() string { return "users" }

func (m userModel) toDomain() domain.User {
	return domain.User{
		ID:            m.ID,
		ExternalID:    m.ExternalID,
		Email:         m.Email,
		WalletAddress: m.WalletAddress,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

func fromDomain(user domain.User) userModel {
	return userModel{
		ID:            user.ID,
		ExternalID:    user.ExternalID,
		Email:         user.Email,
		WalletAddress: user.WalletAddress,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}
}
