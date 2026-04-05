package bootstrap

import (
	"github.com/cashback-platform/services/mint-consumer/internal/infra/database"

	"go.uber.org/fx"
)

var Database = fx.Module("database",
	fx.Provide(database.NewPostgresDB),
)
