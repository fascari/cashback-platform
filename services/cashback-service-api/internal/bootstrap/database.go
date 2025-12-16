package bootstrap

import (
	"github.com/cashback-platform/services/cashback-service-api/internal/config"
	"github.com/cashback-platform/services/cashback-service-api/internal/database"
	"github.com/cashback-platform/services/cashback-service-api/pkg/logger"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

var Database = fx.Module("database",
	fx.Provide(NewDatabase),
)

func NewDatabase(cfg config.Database) (*gorm.DB, error) {
	db, err := database.ConnectPostgres(cfg)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		return nil, err
	}
	return db, nil
}
