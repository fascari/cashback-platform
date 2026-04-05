package database

import (
	"fmt"

	"github.com/cashback-platform/services/mint-consumer/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlg "gorm.io/gorm/logger"
)

func NewPostgresDB(cfg config.Database) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlg.Default.LogMode(gormlg.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	return db, nil
}
