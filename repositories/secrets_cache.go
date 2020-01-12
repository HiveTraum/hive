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
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"strconv"
	"time"
)

func getSecretKey(id models.SecretID) string {
	return fmt.Sprintf("%s:%d", enums.Secret, id)
}

func cacheSecret(cache *redis.Client, span opentracing.Span, secret *models.Secret, key string, timeout time.Duration) error {
	secretCache := &inout.SecretCache{
		Id:      int64(secret.Id),
		Created: secret.Created,
		Value:   secret.Value,
	}

	data, err := proto.Marshal(secretCache)

	if err != nil {
		span.LogFields(log.Error(err))
		sentry.CaptureException(err)
		return nil
	}

	result := cache.Set(key, data, timeout)
	return result.Err()
}

func getSecretFromCache(cache *redis.Client, span opentracing.Span, key string) *models.Secret {

	value, err := cache.Get(key).Bytes()

	if err != nil {
		span.LogFields(log.Error(err))
		sentry.CaptureException(err)
		return nil
	}

	var secretCache inout.SecretCache

	err = proto.Unmarshal(value, &secretCache)

	if err != nil {
		span.LogFields(log.Error(err))
		sentry.CaptureException(err)
		return nil
	}

	span.LogFields(log.String("secret_id", strconv.Itoa(int(secretCache.Id))))

	return &models.Secret{
		Id:      models.SecretID(secretCache.Id),
		Created: secretCache.Created,
		Value:   secretCache.Value,
	}
}

func CacheActualSecret(cache *redis.Client, ctx context.Context, secret *models.Secret, timeout time.Duration) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Cache actual secret")
	err := cacheSecret(cache, span, secret, enums.ActualSecret, timeout)
	span.LogFields(log.String("secret_id", strconv.Itoa(int(secret.Id))))
	span.Finish()
	return err
}

func CacheSecret(cache *redis.Client, ctx context.Context, secret *models.Secret, timeout time.Duration) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Cache secret")
	err := cacheSecret(cache, span, secret, getSecretKey(secret.Id), timeout)
	span.LogFields(log.String("secret_id", strconv.Itoa(int(secret.Id))))
	span.Finish()
	return err
}

func GetActualSecret(cache *redis.Client, ctx context.Context) *models.Secret {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Get actual secret")
	secret := getSecretFromCache(cache, span, enums.ActualSecret)
	span.Finish()
	return secret
}

func GetSecretByID(cache *redis.Client, ctx context.Context, id models.SecretID) *models.Secret {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Get secret")
	secret := getSecretFromCache(cache, span, getSecretKey(id))
	span.Finish()
	return secret
}
