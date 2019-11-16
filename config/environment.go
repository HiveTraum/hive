package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Environment struct {
	DatabaseHost string `env:"DATABASE_HOST" envDefault:"localhost"`
	DatabasePort int    `env:"DATABASE_PORT" envDefault:"5432"`
	DatabaseUser string `env:"DATABASE_USER" envDefault:"auth"`
	DatabasePass string `env:"DATABASE_PASS" envDefault:"123"`
	DatabaseName string `env:"DATABASE_NAME" envDefault:"auth"`
	RedisUrl     string `env:"REDIS_URL" envDefault:"localhost:6379"`
	SentryDsn    string `env:"SENTRY_DSN" envDefault:"https://3e6b6318d35a457dbd57b1445919b38d@sentry.io/1797534"`
}

func GetEnvironment() Environment {
	cfg := Environment{}
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	return cfg
}

func InitEnv() Environment {
	_ = godotenv.Load()
	return GetEnvironment()
}

func isTest() bool {
	return flag.Lookup("test.v") != nil
}
