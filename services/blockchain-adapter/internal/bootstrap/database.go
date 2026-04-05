package bootstrap

import (
	"github.com/cashback-platform/services/blockchain-adapter/internal/infra/database"
	"go.uber.org/fx"
)

var Database = fx.Module("database",
	fx.Provide(database.NewPostgresDB),
)
