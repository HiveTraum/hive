package repositories

import (
	"auth/config"
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCreateEmailConfirmationCode(t *testing.T) {
	cache := config.InitRedis()
	cache.FlushAll()
	code := CreateEmailConfirmationCode(context.Background(), cache, "mail@mail.com", "1234", time.Millisecond)
	require.NotNil(t, code)
	require.Equal(t, "1234", code.Code)
	require.Equal(t, "mail@mail.com", code.Email)
}

func TestGetEmailConfirmationCode(t *testing.T) {
	cache := config.InitRedis()
	cache.FlushAll()
	email := "mail@mail.com"
	CreateEmailConfirmationCode(context.Background(), cache, email, "1234", time.Millisecond)
	code := GetEmailConfirmationCode(context.Background(), cache, email)
	require.Equal(t, "1234", code)
}

func TestGetExpiredEmailConfirmationCode(t *testing.T) {
	cache := config.InitRedis()
	cache.FlushAll()
	email := "mail@mail.com"
	CreateEmailConfirmationCode(context.Background(), cache, email, "1234", time.Millisecond)
	time.Sleep(time.Millisecond * 2)
	code := GetEmailConfirmationCode(context.Background(), cache, email)
	require.Empty(t, code)
}
