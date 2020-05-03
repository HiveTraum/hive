package repositories

import (
	"auth/config"
	"auth/models"
	"context"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetUserViewFromCache(t *testing.T) {
	cache := config.InitRedis(config.InitEnvironment())
	cache.FlushAll()
	ctx := context.Background()
	userView := &models.UserView{
		Id:      uuid.NewV4(),
		Created: 1,
		Roles:   []string{"role"},
		Phones:  []string{"phone"},
		Emails:  []string{"email"},
		RolesID: []uuid.UUID{uuid.NewV4()},
	}
	CacheUserView(cache, ctx, []*models.UserView{userView})
	userViewFromCache := GetUserViewFromCache(cache, ctx, userView.Id)
	require.Equal(t, userView, userViewFromCache)
}

func TestGetUserViewFromCacheWithoutUserView(t *testing.T) {
	cache := config.InitRedis(config.InitEnvironment())
	cache.FlushAll()
	ctx := context.Background()
	userViewFromCache := GetUserViewFromCache(cache, ctx, uuid.NewV4())
	require.Nil(t, userViewFromCache)
}
