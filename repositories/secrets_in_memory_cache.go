package repositories

import (
	"auth/enums"
	"auth/models"
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/patrickmn/go-cache"
	uuid "github.com/satori/go.uuid"
	"time"
)

func cacheSecretInMemory(cache *cache.Cache, secret *models.Secret, key string, timeout time.Duration) {
	cache.Set(key, secret, timeout)
}

func getSecretFromInMemoryCache(cache *cache.Cache, key string) *models.Secret {

	var secret *models.Secret

	if x, found := cache.Get(key); found {
		secret = x.(*models.Secret)
	} else {
		return nil
	}

	return secret
}

func CacheActualSecretInMemory(cache *cache.Cache, ctx context.Context, secret *models.Secret, timeout time.Duration) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cache actual secret in memory")
	cacheSecretInMemory(cache, secret, enums.ActualSecret, timeout)
	span.LogFields(log.String("secret_id", secret.Id.String()))
	span.Finish()
}

func CacheSecretInMemory(cache *cache.Cache, ctx context.Context, secret *models.Secret, timeout time.Duration) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cache secret in memory")
	cacheSecretInMemory(cache, secret, getSecretKey(secret.Id), timeout)
	span.LogFields(log.String("secret_id", secret.Id.String()))
	span.Finish()
}

func GetActualSecretFromInMemoryCache(cache *cache.Cache, ctx context.Context) *models.Secret {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Get actual secret")
	secret := getSecretFromInMemoryCache(cache, enums.ActualSecret)
	span.Finish()
	return secret
}

func GetSecretByIDFromInMemoryCache(cache *cache.Cache, ctx context.Context, id uuid.UUID) *models.Secret {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Get secret")
	secret := getSecretFromInMemoryCache(cache, getSecretKey(id))
	span.Finish()
	return secret
}
