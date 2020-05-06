package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Environment struct {
	Service  string `env:"SERVICE" envDefault:"hive"`
	Instance string `env:"INSTANCE" envDefault:"local"`

	DatabaseURI      string `env:"DATABASE_URI" envDefault:"postgres://hive:123@localhost:5432/hive"`
	RedisURI         string `env:"REDIS_URL" envDefault:"redis://localhost/"`
	SentryDSN        string `env:"SENTRY_DSN"`
	NSQLookupAddress string `env:"NSQ_LOOKUP_ADDRESS" envDefault:"localhost:4180"`

	AccessTokenLifetime    int64  `env:"ACCESS_TOKEN_LIFETIME" envDefault:"15"`  // Minutes
	RefreshTokenLifetime   int64  `env:"REFRESH_TOKEN_LIFETIME" envDefault:"30"` // Days
	RefreshTokenCookieName string `env:"REFRESH_TOKEN_COOKIE_NAME" envDefault:"refreshToken"`
	ActualSecretLifetime   int64  `env:"ACTUAL_SECRET_LIFETIME" envDefault:"1440"` // Minutes
	DefaultPaginationLimit int    `env:"DEFAULT_PAGINATION_LIMIT" envDefault:"50"`
	
	InitialAdmin          string `env:"INITIAL_ADMIN"`
	AdminRole             string `env:"ADMIN_ROLE" envDefault:"admin"`
	RequestContextUserKey string `env:"REQUEST_CONTEXT_USER_KEY" envDefault:"UserContextKey"`
	ServerAddress         string `env:"SERVER_ADDRESS" envDefault:"0.0.0.0:8080"`
	LocalNetworkNamespace string `env:"LOCAL_NETWORK_NAMESPACE" envDefault:"[::1]:"`
}

func InitEnvironment() *Environment {
	_ = godotenv.Load()
	cfg := Environment{}
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}
	log.Log().Msg("Environment successfully parsed")
	return &cfg
}

func isTest() bool {
	return flag.Lookup("test.v") != nil
}
