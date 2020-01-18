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
	uuid "github.com/satori/go.uuid"
	"time"
)

func getSecretKey(id uuid.UUID) string {
	return fmt.Sprintf("%s:%s", enums.Secret, id.String())
}

func cacheSecret(cache *redis.Client, span opentracing.Span, secret *models.Secret, key string, timeout time.Duration) error {
	secretCache := &inout.SecretCache{
		Id:      secret.Id.Bytes(),
		Created: secret.Created,
		Value:   secret.Value.Bytes(),
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

	span.LogFields(log.String("secret_id", uuid.FromBytesOrNil(secretCache.Id).String()))

	return &models.Secret{
		Id:      uuid.FromBytesOrNil(secretCache.Id),
		Created: secretCache.Created,
		Value:   uuid.FromBytesOrNil(secretCache.Value),
	}
}

func CacheActualSecret(cache *redis.Client, ctx context.Context, secret *models.Secret, timeout time.Duration) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Cache actual secret")
	err := cacheSecret(cache, span, secret, enums.ActualSecret, timeout)
	span.LogFields(log.String("secret_id", secret.Id.String()))
	span.Finish()
	return err
}

func CacheSecret(cache *redis.Client, ctx context.Context, secret *models.Secret, timeout time.Duration) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Cache secret")
	err := cacheSecret(cache, span, secret, getSecretKey(secret.Id), timeout)
	span.LogFields(log.String("secret_id", secret.Id.String()))
	span.Finish()
	return err
}

func GetActualSecretFromCache(cache *redis.Client, ctx context.Context) *models.Secret {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Get actual secret")
	secret := getSecretFromCache(cache, span, enums.ActualSecret)
	span.Finish()
	return secret
}

func GetSecretByIDFromCache(cache *redis.Client, ctx context.Context, id uuid.UUID) *models.Secret {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Get secret")
	secret := getSecretFromCache(cache, span, getSecretKey(id))
	span.Finish()
	return secret
}
