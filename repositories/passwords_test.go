package repositories

import (
	"hive/config"
	"hive/enums"
	"context"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreatePassword(t *testing.T) {
	pool := config.InitPool(nil, config.InitEnvironment())
	ctx := context.Background()
	PurgeUsers(pool, ctx)
	PurgePasswords(pool, ctx)
	user := CreateUser(pool, ctx)
	status, password := CreatePassword(pool, ctx, user.Id, "123")
	require.Equal(t, enums.Ok, status)
	require.NotNil(t, password)
}

func TestCreatePasswordWithoutUser(t *testing.T) {
	pool := config.InitPool(nil, config.InitEnvironment())
	ctx := context.Background()
	PurgeUsers(pool, ctx)
	PurgePasswords(pool, ctx)
	status, password := CreatePassword(pool, ctx, uuid.NewV4(), "123")
	require.Equal(t, enums.UserNotFound, status)
	require.Nil(t, password)
}

func TestGetPasswords(t *testing.T) {
	pool := config.InitPool(nil, config.InitEnvironment())
	ctx := context.Background()
	PurgeUsers(pool, ctx)
	PurgePasswords(pool, ctx)
	user := CreateUser(pool, ctx)
	CreatePassword(pool, ctx, user.Id, "123")
	CreatePassword(pool, ctx, user.Id, "456")
	passwords := GetPasswords(pool, ctx, user.Id)
	require.Len(t, passwords, 2)
}

func TestGetLatestPassword(t *testing.T) {
	pool := config.InitPool(nil, config.InitEnvironment())
	ctx := context.Background()
	PurgeUsers(pool, ctx)
	PurgePasswords(pool, ctx)
	user := CreateUser(pool, ctx)
	CreatePassword(pool, ctx, user.Id, "123")
	CreatePassword(pool, ctx, user.Id, "456")
	status, password := GetLatestPassword(pool, ctx, user.Id)
	require.Equal(t, enums.Ok, status)
	require.NotNil(t, password)
	require.Equal(t, "456", password.Value)
	require.Equal(t, user.Id, password.UserId)
}
