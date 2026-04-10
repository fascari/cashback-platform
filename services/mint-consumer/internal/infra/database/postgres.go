package database

import (
	"fmt"

	"github.com/cashback-platform/services/mint-consumer/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlg "gorm.io/gorm/logger"
)

func NewPostgresDB(cfg config.Database) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{
		Logger: gormlg.Default.LogMode(gormlg.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	return db, nil
}
