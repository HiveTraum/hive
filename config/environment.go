package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Environment struct {
	DatabaseHost           string `env:"DATABASE_HOST" envDefault:"localhost"`
	DatabasePort           int    `env:"DATABASE_PORT" envDefault:"5432"`
	DatabaseUser           string `env:"DATABASE_USER" envDefault:"hive"`
	DatabasePass           string `env:"DATABASE_PASS" envDefault:"123"`
	DatabaseName           string `env:"DATABASE_NAME" envDefault:"hive"`
	NSQLookupAddress       string `env:"NSQ_LOOKUP_ADDRESS" envDefault:"localhost:4180"`
	RedisUrl               string `env:"REDIS_URL" envDefault:"localhost:6379"`
	SentryDsn              string `env:"SENTRY_DSN" envDefault:"https://3e6b6318d35a457dbd57b1445919b38d@sentry.io/1797534"`
	EsbUrl                 string `env:"ESB_URL"`
	ESBSender              string `env:"ESB_SENDER" envDefault:"auth"`
	AccessTokenLifetime    int64  `env:"ACCESS_TOKEN_LIFETIME" envDefault:"15"`  // Minutes
	RefreshTokenLifetime   int64  `env:"REFRESH_TOKEN_LIFETIME" envDefault:"30"` // Days
	RefreshTokenCookieName string `env:"REFRESH_TOKEN_COOKIE_NAME" envDefault:"refreshToken"`
	ActualSecretLifetime   int64  `env:"ACTUAL_SECRET_LIFETIME" envDefault:"1440"` // Minutes
	DefaultPaginationLimit int    `env:"DEFAULT_PAGINATION_LIMIT" envDefault:"50"`
	IsTestEnvironment      bool   `env:"IS_TEST_ENVIRONMENT" envDefault:"false"`
	TestConfirmationCode   string `env:"TEST_CONFIRMATION_CODE" envDefault:"111111"`
	InitialAdmin           string `env:"INITIAL_ADMIN"`
	AdminRole              string `env:"ADMIN_ROLE" envDefault:"admin"`
	RequestContextUserKey  string `env:"REQUEST_CONTEXT_USER_KEY" envDefault:"UserContextKey"`
	ServerAddress          string `env:"SERVER_ADDRESS" envDefault:"0.0.0.0:8080"`
	ServiceName            string `env:"SERVICE_NAME" envDefault:"auth"`
	LocalNetworkNamespace  string `env:"LOCAL_NETWORK_NAMESPACE" envDefault:"[::1]:"`
}

func InitEnvironment() *Environment {
	_ = godotenv.Load()
	cfg := Environment{}
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}
	return &cfg
}

func isTest() bool {
	return flag.Lookup("test.v") != nil
}
