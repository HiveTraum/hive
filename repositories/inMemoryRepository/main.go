package inMemoryRepository

import (
	"hive/models"
	"context"
	"github.com/patrickmn/go-cache"
	uuid "github.com/satori/go.uuid"
	"time"
)

type IInMemoryRepository interface {

	// Secrets

	GetSecret(ctx context.Context, id uuid.UUID) *models.Secret
	CacheSecret(ctx context.Context, secret *models.Secret, timeout time.Duration)
	GetActualSecret(ctx context.Context) *models.Secret
	CacheActualSecret(ctx context.Context, secret *models.Secret, timeout time.Duration)
}

type InMemoryRepository struct {
	inMemoryCache *cache.Cache
}

func InitInMemoryRepository(inMemoryCache *cache.Cache) *InMemoryRepository {
	return &InMemoryRepository{
		inMemoryCache: inMemoryCache,
	}
}
