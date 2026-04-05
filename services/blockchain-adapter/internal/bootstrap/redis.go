package bootstrap

import (
	"github.com/cashback-platform/services/blockchain-adapter/internal/config"
	infraredis "github.com/cashback-platform/services/blockchain-adapter/internal/infra/redis"
	"go.uber.org/fx"
)

var Redis = fx.Module("redis",
	fx.Provide(newRedisClient),
)

func newRedisClient(cfg *config.Config) (*infraredis.Client, error) {
	return infraredis.New(cfg.Redis.URL)
}
