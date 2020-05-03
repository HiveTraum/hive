package redisRepository

import (
	"auth/enums"
	"auth/inout"
	"auth/models"
	"context"
	"fmt"
	"github.com/getsentry/sentry-go"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/protobuf/proto"
	"time"
)

func getSecretKey(id uuid.UUID) string {
	return fmt.Sprintf("%s:%s", enums.Secret, id.String())
}

func (repository *RedisRepository) cacheSecret(ctx context.Context, secret *models.Secret, key string, timeout time.Duration) error {
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

	result := repository.redis.WithContext(ctx).Set(key, data, timeout)
	return result.Err()
}

func (repository *RedisRepository) getSecret(ctx context.Context, key string) *models.Secret {

	value, err := repository.redis.WithContext(ctx).Get(key).Bytes()

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

func (repository *RedisRepository) CacheActualSecret(ctx context.Context, secret *models.Secret, timeout time.Duration) error {
	return repository.cacheSecret(ctx, secret, enums.ActualSecret, timeout)
}

func (repository *RedisRepository) CacheSecret(ctx context.Context, secret *models.Secret, timeout time.Duration) error {
	return repository.cacheSecret(ctx, secret, getSecretKey(secret.Id), timeout)
}

func (repository *RedisRepository) GetActualSecret(ctx context.Context) *models.Secret {
	return repository.getSecret(ctx, enums.ActualSecret)
}

func (repository *RedisRepository) GetSecret(ctx context.Context, id uuid.UUID) *models.Secret {
	return repository.getSecret(ctx, getSecretKey(id))
}
