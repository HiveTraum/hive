package config

import (
	"github.com/getsentry/sentry-go"
)

func InitSentry() {

	env := InitEnv()

	err := sentry.Init(sentry.ClientOptions{
		Dsn: env.SentryDsn,
	})

	if err != nil {
		panic(err)
	}
}
