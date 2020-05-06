package config

import (
	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog/log"
)

func InitSentry(environment *Environment) {
	if environment.SentryDSN != "" {
		err := sentry.Init(sentry.ClientOptions{
			Dsn: environment.SentryDSN,
		})

		if err != nil {
			panic(err)
		}

		log.Log().Msg("Sentry successfully initiated")
	} else {
		log.Log().Msg("Provide sentry dsn to instantiate sentry")
	}
}
