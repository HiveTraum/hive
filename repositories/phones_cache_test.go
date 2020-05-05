package repositories

import (
	"hive/config"
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCreatePhoneConfirmationCode(t *testing.T) {
	cache := config.InitRedis(config.InitEnvironment())
	cache.FlushAll()
	ctx := context.Background()
	code := CreatePhoneConfirmationCode(cache, ctx, "+79691234567", "1234", time.Millisecond)
	require.NotNil(t, code)
	require.Equal(t, "1234", code.Code)
	require.Equal(t, "+79691234567", code.Phone)
}

func TestGetPhoneConfirmationCode(t *testing.T) {
	cache := config.InitRedis(config.InitEnvironment())
	cache.FlushAll()
	ctx := context.Background()
	phone := "+79691234567"
	CreatePhoneConfirmationCode(cache, ctx, phone, "1234", time.Millisecond)
	code := GetPhoneConfirmationCode(cache, ctx, phone)
	require.Equal(t, "1234", code)
}

func TestGetExpiredPhoneConfirmationCode(t *testing.T) {
	cache := config.InitRedis(config.InitEnvironment())
	cache.FlushAll()
	ctx := context.Background()
	phone := "+79691234567"
	CreatePhoneConfirmationCode(cache, ctx, phone, "1234", time.Millisecond)
	time.Sleep(time.Millisecond * 2)
	code := GetPhoneConfirmationCode(cache, ctx, phone)
	require.Empty(t, code)
}
