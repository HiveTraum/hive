package repositories

import (
	"auth/models"
	"context"
	"github.com/go-redis/redis/v7"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/patrickmn/go-cache"
	uuid "github.com/satori/go.uuid"
	"time"
)

type IRepository interface {

	// Secrets

	GetSecretByIDFromInMemoryCache(ctx context.Context, id uuid.UUID) *models.Secret
	GetSecretByIDFromCache(ctx context.Context, id uuid.UUID) *models.Secret
	CacheSecretInMemory(ctx context.Context, secret *models.Secret, timeout time.Duration)
	GetSecretFromDB(ctx context.Context, id uuid.UUID) *models.Secret
	CacheSecret(ctx context.Context, secret *models.Secret, timeout time.Duration) error
}

type Repository struct {
	postgres      *pgxpool.Pool
	redis         *redis.Client
	inMemoryCache *cache.Cache
}

func InitRepository(postgres *pgxpool.Pool, redis *redis.Client, inMemoryCache *cache.Cache) *Repository {
	return &Repository{
		postgres:      postgres,
		redis:         redis,
		inMemoryCache: inMemoryCache,
	}
}
