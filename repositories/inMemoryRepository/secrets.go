package inMemoryRepository

import (
	"auth/enums"
	"auth/models"
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	uuid "github.com/satori/go.uuid"
	"time"
)

func (repository *InMemoryRepository) cacheSecret(secret *models.Secret, key string, timeout time.Duration) {
	repository.inMemoryCache.Set(key, secret, timeout)
}

func (repository *InMemoryRepository) getSecret(key string) *models.Secret {

	var secret *models.Secret
	if x, found := repository.inMemoryCache.Get(key); found {
		secret = x.(*models.Secret)
	} else {
		return nil
	}
	return secret
}

func getSecretKey(id uuid.UUID) string {
	return fmt.Sprintf("%s:%s", enums.Secret, id.String())
}

func (repository *InMemoryRepository) CacheActualSecret(ctx context.Context, secret *models.Secret, timeout time.Duration) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cache actual secret in memory")
	repository.cacheSecret(secret, enums.ActualSecret, timeout)
	span.LogFields(log.String("secret_id", secret.Id.String()))
	span.Finish()
}

func (repository *InMemoryRepository) CacheSecret(ctx context.Context, secret *models.Secret, timeout time.Duration) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cache secret in memory")
	repository.cacheSecret(secret, getSecretKey(secret.Id), timeout)
	span.LogFields(log.String("secret_id", secret.Id.String()))
	span.Finish()
}

func (repository *InMemoryRepository) GetActualSecret(ctx context.Context) *models.Secret {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Get actual secret")
	secret := repository.getSecret(enums.ActualSecret)
	span.Finish()
	return secret
}

func (repository *InMemoryRepository) GetSecret(ctx context.Context, id uuid.UUID) *models.Secret {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Get secret")
	secret := repository.getSecret(getSecretKey(id))
	span.Finish()
	return secret
}
