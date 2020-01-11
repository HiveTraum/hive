package repositories

import (
	"auth/config"
	"auth/enums"
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateEmail(t *testing.T) {
	pool := config.InitPool(nil)
	ctx := context.Background()
	PurgeUsers(pool, ctx)
	PurgeEmails(pool, ctx)
	user := CreateUser(pool, ctx)
	status, email := CreateEmail(pool, ctx, user.Id, "mail@mail.com")
	require.Equal(t, enums.Ok, status)
	require.NotNil(t, email)
}

func TestCreateEmailWithoutUser(t *testing.T) {
	pool := config.InitPool(nil)
	ctx := context.Background()
	PurgeUsers(pool, ctx)
	PurgeEmails(pool, ctx)
	status, email := CreateEmail(pool, ctx, 1, "mail@mail.com")
	require.Equal(t, enums.UserNotFound, status)
	require.Nil(t, email)
}

func TestGetEmail(t *testing.T) {
	pool := config.InitPool(nil)
	ctx := context.Background()
	PurgeUsers(pool, ctx)
	PurgeEmails(pool, ctx)
	user := CreateUser(pool, ctx)
	CreateEmail(pool, ctx, user.Id, "mail@mail.com")
	status, email := GetEmail(pool, ctx, "mail@mail.com")
	require.Equal(t, enums.Ok, status)
	require.NotNil(t, email)
	require.Equal(t, "mail@mail.com", email.Value)
}
