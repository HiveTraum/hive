package repositories

import (
	"auth/config"
	"auth/models"
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCacheSecret(t *testing.T) {
	cache := config.InitRedis()
	cache.FlushAll()
	ctx := context.Background()
	secret := &models.Secret{
		Id:      1,
		Created: 1,
		Value:   "Hello",
	}
	err := CacheSecret(cache, ctx, secret, time.Millisecond)
	require.Nil(t, err)
}

func TestGetSecretByIDSecret(t *testing.T) {
	cache := config.InitRedis()
	cache.FlushAll()
	ctx := context.Background()

	secret := &models.Secret{
		Id:      1,
		Created: 1,
		Value:   "Hello",
	}

	_ = CacheSecret(cache, ctx, secret, time.Millisecond)

	cachedSecret := GetSecretByID(cache, ctx, 1)
	require.NotNil(t, cachedSecret)
	require.Equal(t, secret, cachedSecret)
}

func TestGetExpiredSecretByIDSecret(t *testing.T) {
	cache := config.InitRedis()
	cache.FlushAll()
	ctx := context.Background()

	secret := &models.Secret{
		Id:      1,
		Created: 1,
		Value:   "Hello",
	}

	_ = CacheSecret(cache, ctx, secret, time.Millisecond)
	time.Sleep(time.Millisecond * 2)
	cachedSecret := GetSecretByID(cache, ctx, 1)
	require.Nil(t, cachedSecret)
}

func TestCacheActualSecret(t *testing.T) {
	cache := config.InitRedis()
	cache.FlushAll()
	ctx := context.Background()
	secret := &models.Secret{
		Id:      1,
		Created: 1,
		Value:   "Hello",
	}
	err := CacheActualSecret(cache, ctx, secret, time.Millisecond)
	require.Nil(t, err)
}

func TestGetActualSecret(t *testing.T) {
	cache := config.InitRedis()
	cache.FlushAll()
	ctx := context.Background()
	secret := &models.Secret{
		Id:      1,
		Created: 1,
		Value:   "Hello",
	}
	_ = CacheActualSecret(cache, ctx, secret, time.Millisecond)
	cachedSecret := GetActualSecret(cache, ctx)
	require.NotNil(t, cachedSecret)
	require.Equal(t, secret, cachedSecret)
}
