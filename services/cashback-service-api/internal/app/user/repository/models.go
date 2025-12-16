package repository

import (
	"time"

	"github.com/cashback-platform/services/cashback-service-api/internal/app/user/domain"
	"github.com/google/uuid"
)

// userModel represents the database model for users
type userModel struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ExternalID    string    `gorm:"uniqueIndex;not null"`
	Email         string    `gorm:"uniqueIndex;not null"`
	WalletAddress string    `gorm:"not null"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

func (userModel) TableName() string {
	return "users"
}

// toDomain converts database model to domain entity
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

// fromDomain converts domain entity to database model
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
