package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"sync"
)

type Environment struct {
	DatabaseHost           string `env:"DATABASE_HOST" envDefault:"localhost"`
	DatabasePort           int    `env:"DATABASE_PORT" envDefault:"5432"`
	DatabaseUser           string `env:"DATABASE_USER" envDefault:"auth"`
	DatabasePass           string `env:"DATABASE_PASS" envDefault:"123"`
	DatabaseName           string `env:"DATABASE_NAME" envDefault:"auth"`
	RedisUrl               string `env:"REDIS_URL" envDefault:"localhost:6379"`
	SentryDsn              string `env:"SENTRY_DSN" envDefault:"https://3e6b6318d35a457dbd57b1445919b38d@sentry.io/1797534"`
	EsbUrl                 string `env:"ESB_URL"`
	AccessTokenLifetime    int64  `env:"ACCESS_TOKEN_LIFETIME" envDefault:"15"`    // Minutes
	RefreshTokenLifetime   int64  `env:"REFRESH_TOKEN_LIFETIME" envDefault:"30"`   // Days
	ActualSecretLifetime   int64  `env:"ACTUAL_SECRET_LIFETIME" envDefault:"1440"` // Minutes
	DefaultPaginationLimit int    `env:"DEFAULT_PAGINATION_LIMIT" envDefault:"50"`
}

var cfg Environment
var onceEnvironment sync.Once

func GetEnvironment() Environment {
	onceEnvironment.Do(func() {
		cfg = Environment{}
		if err := env.Parse(&cfg); err != nil {
			panic(err)
		}
	})

	return cfg
}

func InitEnv() Environment {
	_ = godotenv.Load()
	return GetEnvironment()
}

func isTest() bool {
	return flag.Lookup("test.v") != nil
}
