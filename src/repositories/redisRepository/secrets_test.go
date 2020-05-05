package redisRepository

import (
	"hive/config"
	"hive/models"
	"context"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCacheSecret(t *testing.T) {
	cache := config.InitRedis(config.InitEnvironment())
	repo := InitRedisRepository(cache)
	cache.FlushAll()
	ctx := context.Background()
	secret := &models.Secret{
		Id:      uuid.NewV4(),
		Created: 1,
		Value:   uuid.NewV4(),
	}
	err := repo.CacheSecret(ctx, secret, time.Millisecond)
	require.Nil(t, err)
}

func TestGetSecretByIDSecret(t *testing.T) {
	cache := config.InitRedis(config.InitEnvironment())
	repo := InitRedisRepository(cache)
	cache.FlushAll()
	ctx := context.Background()

	id := uuid.NewV4()

	secret := &models.Secret{
		Id:      id,
		Created: 1,
		Value:   uuid.NewV4(),
	}

	_ = repo.CacheSecret(ctx, secret, time.Millisecond)

	cachedSecret := repo.GetSecret(ctx, id)
	require.NotNil(t, cachedSecret)
	require.Equal(t, secret, cachedSecret)
}

func TestGetExpiredSecretByIDSecret(t *testing.T) {
	cache := config.InitRedis(config.InitEnvironment())
	repo := InitRedisRepository(cache)
	cache.FlushAll()
	ctx := context.Background()

	id := uuid.NewV4()

	secret := &models.Secret{
		Id:      id,
		Created: 1,
		Value:   uuid.NewV4(),
	}

	_ = repo.CacheSecret(ctx, secret, time.Millisecond)
	time.Sleep(time.Millisecond * 2)
	cachedSecret := repo.GetSecret(ctx, id)
	require.Nil(t, cachedSecret)
}

func TestCacheActualSecret(t *testing.T) {
	cache := config.InitRedis(config.InitEnvironment())
	repo := InitRedisRepository(cache)
	cache.FlushAll()
	ctx := context.Background()

	secret := &models.Secret{
		Id:      uuid.NewV4(),
		Created: 1,
		Value:   uuid.NewV4(),
	}
	err := repo.CacheActualSecret(ctx, secret, time.Millisecond)
	require.Nil(t, err)
}

func TestGetActualSecret(t *testing.T) {
	cache := config.InitRedis(config.InitEnvironment())
	repo := InitRedisRepository(cache)
	cache.FlushAll()
	ctx := context.Background()
	secret := &models.Secret{
		Id:      uuid.NewV4(),
		Created: 1,
		Value:   uuid.NewV4(),
	}
	_ = repo.CacheActualSecret(ctx, secret, time.Millisecond)
	cachedSecret := repo.GetActualSecret(ctx)
	require.NotNil(t, cachedSecret)
	require.Equal(t, secret, cachedSecret)
}
