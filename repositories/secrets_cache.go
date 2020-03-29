package repositories

import (
	"auth/enums"
	"auth/inout"
	"auth/models"
	"context"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/go-redis/redis/v7"
	"github.com/golang/protobuf/proto"
	uuid "github.com/satori/go.uuid"
	"time"
)

func getSecretKey(id uuid.UUID) string {
	return fmt.Sprintf("%s:%s", enums.Secret, id.String())
}

func cacheSecret(cache *redis.Client, ctx context.Context, secret *models.Secret, key string, timeout time.Duration) error {
	secretCache := &inout.SecretCache{
		Id:      secret.Id.Bytes(),
		Created: secret.Created,
		Value:   secret.Value.Bytes(),
	}

	data, err := proto.Marshal(secretCache)

	if err != nil {
		sentry.CaptureException(err)
		return nil
	}

	result := cache.WithContext(ctx).Set(key, data, timeout)
	return result.Err()
}

func getSecretFromCache(ctx context.Context, cache *redis.Client, key string) *models.Secret {

	value, err := cache.WithContext(ctx).Get(key).Bytes()

	if err != nil {
		return nil
	}

	var secretCache inout.SecretCache

	err = proto.Unmarshal(value, &secretCache)

	if err != nil {
		sentry.CaptureException(err)
		return nil
	}

	return &models.Secret{
		Id:      uuid.FromBytesOrNil(secretCache.Id),
		Created: secretCache.Created,
		Value:   uuid.FromBytesOrNil(secretCache.Value),
	}
}

func CacheActualSecret(cache *redis.Client, ctx context.Context, secret *models.Secret, timeout time.Duration) error {
	return cacheSecret(cache, ctx, secret, enums.ActualSecret, timeout)
}

func CacheSecret(cache *redis.Client, ctx context.Context, secret *models.Secret, timeout time.Duration) error {
	return cacheSecret(cache, ctx, secret, getSecretKey(secret.Id), timeout)
}

func GetActualSecretFromCache(cache *redis.Client, ctx context.Context) *models.Secret {
	return getSecretFromCache(ctx, cache, enums.ActualSecret)
}

func GetSecretByIDFromCache(cache *redis.Client, ctx context.Context, id uuid.UUID) *models.Secret {
	return getSecretFromCache(ctx, cache, getSecretKey(id))
}
