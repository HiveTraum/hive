package repositories

import (
	"auth/config"
	"auth/models"
	"context"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGetActualSecretFromInMemoryCache(t *testing.T) {
	cache := config.InitInMemoryCache()
	cache.Flush()
	ctx := context.Background()
	secret := &models.Secret{
		Id:      uuid.NewV4(),
		Created: 1,
		Value:   uuid.NewV4(),
	}
	CacheActualSecretInMemory(cache, ctx, secret, time.Millisecond)
	secretFromCache := GetActualSecretFromInMemoryCache(cache, ctx)
	require.NotNil(t, secretFromCache)
	require.Equal(t, secret, secretFromCache)
}

func TestGetExpiredActualSecretFromInMemoryCache(t *testing.T) {
	cache := config.InitInMemoryCache()
	cache.Flush()
	ctx := context.Background()
	secret := &models.Secret{
		Id:      uuid.NewV4(),
		Created: 1,
		Value:   uuid.NewV4(),
	}
	CacheActualSecretInMemory(cache, ctx, secret, time.Millisecond)
	time.Sleep(time.Millisecond * 2)
	secretFromCache := GetActualSecretFromInMemoryCache(cache, ctx)
	require.Nil(t, secretFromCache)
}

func TestGetSecretByIDFromInMemoryCache(t *testing.T) {
	cache := config.InitInMemoryCache()
	cache.Flush()
	ctx := context.Background()
	secret := &models.Secret{
		Id:      uuid.NewV4(),
		Created: 1,
		Value:   uuid.NewV4(),
	}
	CacheSecretInMemory(cache, ctx, secret, time.Millisecond)
	secretFromCache := GetSecretByIDFromInMemoryCache(cache, ctx, secret.Id)
	require.NotNil(t, secretFromCache)
	require.Equal(t, secret, secretFromCache)
}

func TestGetExpiredSecretByIDFromInMemoryCache(t *testing.T) {
	cache := config.InitInMemoryCache()
	cache.Flush()
	ctx := context.Background()
	secret := &models.Secret{
		Id:      uuid.NewV4(),
		Created: 1,
		Value:   uuid.NewV4(),
	}
	CacheSecretInMemory(cache, ctx, secret, time.Millisecond)
	time.Sleep(time.Millisecond * 2)
	secretFromCache := GetSecretByIDFromInMemoryCache(cache, ctx, secret.Id)
	require.Nil(t, secretFromCache)
}