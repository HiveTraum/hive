package stores

import (
	"auth/config"
	"auth/models"
	"auth/repositories"
	"context"
	"time"
)

func (store *DatabaseStore) GetSecret(ctx context.Context, id models.SecretID) *models.Secret {

	secret := repositories.GetSecretByID(store.Cache, ctx, id)
	if secret != nil {
		return secret
	}

	secret = repositories.GetSecretFromDB(store.Db, ctx, id)
	repositories.CacheSecret(store.Cache, ctx, secret, time.Hour*48)
	return secret
}

func (store *DatabaseStore) GetActualSecret(ctx context.Context) *models.Secret {

	env := config.GetEnvironment()

	actualSecret := repositories.GetActualSecret(store.Cache, ctx)

	if actualSecret != nil {
		return actualSecret
	}

	actualSecret = repositories.CreateSecret(store.Db, ctx)
	repositories.CacheActualSecret(store.Cache, ctx, actualSecret, time.Minute*env.ActualSecretLifetime)
	repositories.CacheSecret(store.Cache, ctx, actualSecret, time.Hour*env.RefreshTokenLifetime*24)
	return actualSecret
}
