package repositories

import (
	"auth/config"
	"auth/enums"
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateOrUpdateAllUsersViewOnUserCreation(t *testing.T) {
	pool := config.InitPool(nil)
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
	pool := config.InitPool(nil)
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

func TestCreateOrUpdateUsersViewWithTheSamePhone(t *testing.T) {
	pool := config.InitPool(nil)
	ctx := context.Background()
	PurgeUserViews(pool, ctx)
	PurgeUsers(pool, ctx)
	PurgePhones(pool, ctx)
	firstUser := CreateUser(pool, ctx)
	status, phone := CreatePhone(pool, ctx, firstUser.Id, "+71234567890")
	require.Equal(t, enums.Ok, status)
	views := CreateOrUpdateUsersView(pool, ctx, CreateOrUpdateUsersViewQuery{GetUsersQuery: GetUsersQuery{Id: []int64{firstUser.Id}}})
	require.Len(t, views, 1)
	require.Len(t, views[0].Phones, 1)
	require.Equal(t, phone.Value, views[0].Phones[0])
	secondsUserWithTheSamePhone := CreateUser(pool, ctx)
	status, phone = CreatePhone(pool, ctx, secondsUserWithTheSamePhone.Id, phone.Value)
	require.Equal(t, enums.Ok, status)
	require.Equal(t, secondsUserWithTheSamePhone.Id, phone.UserId)
	views = CreateOrUpdateUsersView(pool, ctx, CreateOrUpdateUsersViewQuery{
		GetUsersQuery: GetUsersQuery{
			Id: []int64{firstUser.Id, secondsUserWithTheSamePhone.Id},
		},
	})
	require.Len(t, views, 2)
	require.Len(t, views[0].Phones, 0)
	require.Len(t, views[1].Phones, 1)
}
