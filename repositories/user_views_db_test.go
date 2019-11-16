package repositories

import (
	"auth/config"
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateOrUpdateAllUsersViewOnUserCreation(t *testing.T) {
	pool := config.InitPool()
	ctx := context.Background()
	PurgeUsers(pool, ctx)
	PurgeUserViews(pool, ctx)

	views := CreateOrUpdateUsersView(pool, ctx, CreateOrUpdateUsersViewQuery{})
	require.Len(t, views, 0)
	CreateUser(pool, ctx)
	views = CreateOrUpdateUsersView(pool, ctx, CreateOrUpdateUsersViewQuery{})
	require.Len(t, views, 1)
}

func TestCreateOrUpdateUsersViewOnUserCreation(t *testing.T) {
	pool := config.InitPool()
	ctx := context.Background()
	PurgeUsers(pool, ctx)
	PurgeUserViews(pool, ctx)
	views := CreateOrUpdateUsersView(pool, ctx, CreateOrUpdateUsersViewQuery{})
	require.Len(t, views, 0)
	user := CreateUser(pool, ctx)
	CreateUser(pool, ctx)
	CreateUser(pool, ctx)
	CreateUser(pool, ctx)
	CreateUser(pool, ctx)
	views = CreateOrUpdateUsersView(pool, ctx, CreateOrUpdateUsersViewQuery{
		GetUsersQuery: GetUsersQuery{Id: []int64{user.Id}}, Roles: nil,
	})
	require.Len(t, views, 1)
	require.Equal(t, user.Id, views[0].Id)

	views = CreateOrUpdateUsersView(pool, ctx, CreateOrUpdateUsersViewQuery{})
	require.Len(t, views, 4)
}
