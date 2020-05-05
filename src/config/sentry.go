package config

import (
	"github.com/getsentry/sentry-go"
)

func InitSentry(environment *Environment) {

	err := sentry.Init(sentry.ClientOptions{
		Dsn: environment.SentryDsn,
	})

	if err != nil {
		panic(err)
	}
}
