package stores

import (
	"auth/models"
	"auth/repositories"
	"context"
	"github.com/getsentry/sentry-go"
	uuid "github.com/satori/go.uuid"
	"time"
)

func (store *DatabaseStore) GetSecret(ctx context.Context, id uuid.UUID) *models.Secret {

	secret := repositories.GetSecretByIDFromInMemoryCache(store.inMemoryCache, ctx, id)
	if secret != nil {
		return secret
	}

	secret = repositories.GetSecretByIDFromCache(store.cache, ctx, id)
	if secret != nil {
		repositories.CacheSecretInMemory(store.inMemoryCache, ctx, secret, time.Hour*24)
		return secret
	}

	secret = repositories.GetSecretFromDB(store.db, ctx, id)
	err := repositories.CacheSecret(store.cache, ctx, secret, time.Hour*48)
	if err != nil {
		sentry.CaptureException(err)
	}
	repositories.CacheSecretInMemory(store.inMemoryCache, ctx, secret, time.Hour*24)
	return secret
}

func (store *DatabaseStore) GetActualSecret(ctx context.Context) *models.Secret {

	env := store.environment

	actualSecret := repositories.GetActualSecretFromInMemoryCache(store.inMemoryCache, ctx)
	if actualSecret != nil {
		return actualSecret
	}

	actualSecret = repositories.GetActualSecretFromCache(store.cache, ctx)
	if actualSecret != nil {
		repositories.CacheActualSecretInMemory(store.inMemoryCache, ctx, actualSecret, time.Minute*time.Duration(env.ActualSecretLifetime))
		return actualSecret
	}

	return nil
}

func (store *DatabaseStore) CreateSecret(ctx context.Context) *models.Secret {

	env := store.environment

	actualSecret := repositories.CreateSecret(store.db, ctx)
	err := repositories.CacheActualSecret(store.cache, ctx, actualSecret, time.Minute*time.Duration(env.ActualSecretLifetime))
	if err != nil {
		sentry.CaptureException(err)
	}
	err = repositories.CacheSecret(store.cache, ctx, actualSecret, time.Hour*time.Duration(env.RefreshTokenLifetime)*24)
	if err != nil {
		sentry.CaptureException(err)
	}
	repositories.CacheActualSecretInMemory(store.inMemoryCache, ctx, actualSecret, time.Minute*time.Duration(env.ActualSecretLifetime))
	repositories.CacheSecretInMemory(store.inMemoryCache, ctx, actualSecret, time.Hour*time.Duration(env.RefreshTokenLifetime))
	return actualSecret
}
