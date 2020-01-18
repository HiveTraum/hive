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

	secret := repositories.GetSecretByIDFromInMemoryCache(store.InMemoryCache, ctx, id)
	if secret != nil {
		return secret
	}

	secret = repositories.GetSecretByIDFromCache(store.Cache, ctx, id)
	if secret != nil {
		repositories.CacheSecretInMemory(store.InMemoryCache, ctx, secret, time.Hour*24)
		return secret
	}

	secret = repositories.GetSecretFromDB(store.Db, ctx, id)
	err := repositories.CacheSecret(store.Cache, ctx, secret, time.Hour*48)
	if err != nil {
		sentry.CaptureException(err)
	}
	repositories.CacheSecretInMemory(store.InMemoryCache, ctx, secret, time.Hour*24)
	return secret
}

func (store *DatabaseStore) GetActualSecret(ctx context.Context) *models.Secret {

	env := config.GetEnvironment()

	actualSecret := repositories.GetActualSecretFromInMemoryCache(store.InMemoryCache, ctx)
	if actualSecret != nil {
		return actualSecret
	}

	actualSecret = repositories.GetActualSecretFromCache(store.Cache, ctx)
	if actualSecret != nil {
		repositories.CacheActualSecretInMemory(store.InMemoryCache, ctx, actualSecret, time.Minute*time.Duration(env.ActualSecretLifetime))
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
	repositories.CacheActualSecretInMemory(store.InMemoryCache, ctx, actualSecret, time.Minute*time.Duration(env.ActualSecretLifetime))
	repositories.CacheSecretInMemory(store.InMemoryCache, ctx, actualSecret, time.Hour*time.Duration(env.RefreshTokenLifetime))
	return actualSecret
}
