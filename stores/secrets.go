package stores

import (
	"auth/config"
	"auth/models"
	"auth/repositories"
	"context"
	"github.com/getsentry/sentry-go"
	uuid "github.com/satori/go.uuid"
	"time"
)

func (store *DatabaseStore) GetSecret(ctx context.Context, id uuid.UUID) *models.Secret {

	secret := repositories.GetSecretByID(store.Cache, ctx, id)
	if secret != nil {
		return secret
	}

	secret = repositories.GetSecretFromDB(store.Db, ctx, id)
	err := repositories.CacheSecret(store.Cache, ctx, secret, time.Hour*48)
	if err != nil {
		sentry.CaptureException(err)
	}
	return secret
}

func (store *DatabaseStore) GetActualSecret(ctx context.Context) *models.Secret {

	env := config.GetEnvironment()

	actualSecret := repositories.GetActualSecret(store.Cache, ctx)

	if actualSecret != nil {
		return actualSecret
	}

	actualSecret = repositories.CreateSecret(store.Db, ctx)
	err := repositories.CacheActualSecret(store.Cache, ctx, actualSecret, time.Minute*time.Duration(env.ActualSecretLifetime))
	if err != nil {
		sentry.CaptureException(err)
	}
	err = repositories.CacheSecret(store.Cache, ctx, actualSecret, time.Hour*time.Duration(env.RefreshTokenLifetime)*24)
	if err != nil {
		sentry.CaptureException(err)
	}
	return actualSecret
}
