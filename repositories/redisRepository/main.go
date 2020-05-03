package redisRepository

import (
	"auth/models"
	"context"
	"github.com/go-redis/redis/v7"
	uuid "github.com/satori/go.uuid"
	"time"
)

type IRedisRepository interface {

	// Secrets

	GetSecret(ctx context.Context, id uuid.UUID) *models.Secret
	CacheSecret(ctx context.Context, secret *models.Secret, timeout time.Duration) error
	GetActualSecret(ctx context.Context) *models.Secret
	CacheActualSecret(ctx context.Context, secret *models.Secret, timeout time.Duration) error
}

type RedisRepository struct {
	redis *redis.Client
}

func InitRedisRepository(redis *redis.Client) *RedisRepository {
	return &RedisRepository{
		redis: redis,
	}
}
