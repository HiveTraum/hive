package repositories

import (
	"auth/config"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCreateEmailConfirmationCode(t *testing.T) {
	cache := config.InitRedis()
	cache.FlushAll()
	code := CreateEmailConfirmationCode(cache, "mail@mail.com", "1234", time.Millisecond)
	require.NotNil(t, code)
	require.Equal(t, "1234", code.Code)
	require.Equal(t, "mail@mail.com", code.Email)
}

func TestGetEmailConfirmationCode(t *testing.T) {
	cache := config.InitRedis()
	cache.FlushAll()
	email := "mail@mail.com"
	CreateEmailConfirmationCode(cache, email, "1234", time.Millisecond)
	code := GetEmailConfirmationCode(cache, email)
	require.Equal(t, "1234", code)
}

func TestGetExpiredEmailConfirmationCode(t *testing.T) {
	cache := config.InitRedis()
	cache.FlushAll()
	email := "mail@mail.com"
	CreateEmailConfirmationCode(cache, email, "1234", time.Millisecond)
	time.Sleep(time.Millisecond * 2)
	code := GetEmailConfirmationCode(cache, email)
	require.Empty(t, code)
}
