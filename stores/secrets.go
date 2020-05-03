package stores

import (
	"auth/models"
	"context"
	"github.com/getsentry/sentry-go"
	uuid "github.com/satori/go.uuid"
	"time"
)

func (store *DatabaseStore) GetSecret(ctx context.Context, id uuid.UUID) *models.Secret {

	secret := store.inMemoryRepository.GetSecret(ctx, id)
	if secret != nil {
		return secret
	}

	secret = store.redisRepository.GetSecret(ctx, id)
	if secret != nil {
		store.inMemoryRepository.CacheSecret(ctx, secret, time.Hour*24)
		return secret
	}

	secret = store.postgresRepository.GetSecret(ctx, id)
	if secret != nil {
		err := store.redisRepository.CacheSecret(ctx, secret, time.Hour*48)
		if err != nil {
			sentry.CaptureException(err)
		}
		store.inMemoryRepository.CacheSecret(ctx, secret, time.Hour*24)
	}

	return secret
}

func (store *DatabaseStore) GetActualSecret(ctx context.Context) *models.Secret {

	env := store.environment

	actualSecret := store.inMemoryRepository.GetActualSecret(ctx)
	if actualSecret != nil {
		return actualSecret
	}

	actualSecret = store.redisRepository.GetActualSecret(ctx)
	if actualSecret != nil {
		store.inMemoryRepository.CacheActualSecret(ctx, actualSecret, time.Minute*time.Duration(env.ActualSecretLifetime))
		return actualSecret
	}

	return nil
}

func (store *DatabaseStore) CreateSecret(ctx context.Context) *models.Secret {

	env := store.environment

	actualSecret := store.postgresRepository.CreateSecret(ctx)
	err := store.redisRepository.CacheActualSecret(ctx, actualSecret, time.Minute*time.Duration(env.ActualSecretLifetime))
	if err != nil {
		sentry.CaptureException(err)
	}
	err = store.redisRepository.CacheSecret(ctx, actualSecret, time.Hour*time.Duration(env.RefreshTokenLifetime)*24)
	if err != nil {
		sentry.CaptureException(err)
	}
	store.inMemoryRepository.CacheActualSecret(ctx, actualSecret, time.Minute*time.Duration(env.ActualSecretLifetime))
	store.inMemoryRepository.CacheSecret(ctx, actualSecret, time.Hour*time.Duration(env.RefreshTokenLifetime))
	return actualSecret
}
