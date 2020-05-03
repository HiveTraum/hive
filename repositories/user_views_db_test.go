package repositories

import (
	"auth/config"
	"auth/enums"
	"auth/functools"
	"context"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCreateOrUpdateAllUsersViewOnUserCreation(t *testing.T) {
	pool := config.InitPool(nil, config.InitEnvironment())
	ctx := context.Background()
	PurgeUsers(pool, ctx)
	PurgeUserViews(pool, ctx)

	views := CreateOrUpdateUsersView(pool, ctx, CreateOrUpdateUsersViewStoreQuery{})
	require.Len(t, views, 0)
	CreateUser(pool, ctx)
	views = CreateOrUpdateUsersView(pool, ctx, CreateOrUpdateUsersViewStoreQuery{})
	require.Len(t, views, 1)
}

func TestCreateOrUpdateUsersViewOnUserCreation(t *testing.T) {
	pool := config.InitPool(nil, config.InitEnvironment())
	ctx := context.Background()
	PurgeUsers(pool, ctx)
	PurgeUserViews(pool, ctx)
	views := CreateOrUpdateUsersView(pool, ctx, CreateOrUpdateUsersViewStoreQuery{})
	require.Len(t, views, 0)
	user := CreateUser(pool, ctx)
	CreateUser(pool, ctx)
	CreateUser(pool, ctx)
	CreateUser(pool, ctx)
	CreateUser(pool, ctx)
	views = CreateOrUpdateUsersView(pool, ctx, CreateOrUpdateUsersViewStoreQuery{
		Id: []uuid.UUID{user.Id}, Roles: nil,
		Limit: 10,
	})
	require.Len(t, views, 1)
	require.Equal(t, user.Id, views[0].Id)

	views = CreateOrUpdateUsersView(pool, ctx, CreateOrUpdateUsersViewStoreQuery{})
	require.Len(t, views, 4)
}

func TestCreateOrUpdateUsersViewWithTheSamePhone(t *testing.T) {
	pool := config.InitPool(nil, config.InitEnvironment())
	ctx := context.Background()
	PurgeUserViews(pool, ctx)
	PurgeUsers(pool, ctx)
	PurgePhones(pool, ctx)
	phoneValue := functools.NormalizePhone("+79234567890")
	firstUser := CreateUser(pool, ctx)
	status, phone := CreatePhone(pool, ctx, firstUser.Id, phoneValue)
	require.Equal(t, enums.Ok, status)
	views := CreateOrUpdateUsersView(pool, ctx, CreateOrUpdateUsersViewStoreQuery{Id: []uuid.UUID{firstUser.Id}})
	require.Len(t, views, 1)
	require.Len(t, views[0].Phones, 1)
	require.Equal(t, phone.Value, views[0].Phones[0])
	secondsUserWithTheSamePhone := CreateUser(pool, ctx)
	status, phone = CreatePhone(pool, ctx, secondsUserWithTheSamePhone.Id, phone.Value)
	require.Equal(t, enums.Ok, status)
	require.Equal(t, secondsUserWithTheSamePhone.Id, phone.UserId)
	views = CreateOrUpdateUsersView(pool, ctx, CreateOrUpdateUsersViewStoreQuery{
		Id: []uuid.UUID{firstUser.Id, secondsUserWithTheSamePhone.Id},
	})
	require.Len(t, views, 2)
	require.Len(t, views[0].Phones, 0)
	require.Len(t, views[1].Phones, 1)
}

func TestDatabaseStore_GetUsersViewPagination(t *testing.T) {
	pool := config.InitPool(nil, config.InitEnvironment())
	ctx := context.Background()
	config.PurgeUsers(pool, ctx)
	config.PurgeUserViews(pool, ctx)

	CreateUser(pool, ctx)
	CreateUser(pool, ctx)
	CreateOrUpdateUsersView(pool, ctx, CreateOrUpdateUsersViewStoreQuery{})

	userViews, pagination := GetUsersView(pool, ctx, GetUsersViewStoreQuery{
		Limit: 2,
		Page:  0,
	})

	require.Len(t, userViews, 2)
	require.Equal(t, int64(2), pagination.Count)
	require.False(t, pagination.HasNext)
	require.False(t, pagination.HasPrevious)
}

func TestDatabaseStore_GetUsersViewPaginationWithLimit(t *testing.T) {
	pool := config.InitPool(nil, config.InitEnvironment())
	ctx := context.Background()
	config.PurgeUsers(pool, ctx)
	config.PurgeUserViews(pool, ctx)

	user := CreateUser(pool, ctx)
	CreateUser(pool, ctx)
	CreateOrUpdateUsersView(pool, ctx, CreateOrUpdateUsersViewStoreQuery{})

	userViews, pagination := GetUsersView(pool, ctx, GetUsersViewStoreQuery{
		Limit: 1,
		Page:  1,
	})

	require.Len(t, userViews, 1)
	require.Equal(t, int64(2), pagination.Count)
	require.True(t, pagination.HasNext)
	require.False(t, pagination.HasPrevious)
	require.Equal(t, user.Id, userViews[0].Id)
}

func TestDatabaseStore_GetUsersViewPaginationWithLimitWithoutNext(t *testing.T) {
	pool := config.InitPool(nil, config.InitEnvironment())
	ctx := context.Background()
	config.PurgeUsers(pool, ctx)
	config.PurgeUserViews(pool, ctx)

	CreateUser(pool, ctx)

	// TODO для стабильности тестов ожидаем миллисекунду, т.к. сортировка идет по полю created, которое иногда может быть одинаковым ввиду того что пользователи создаются почти в один момент
	time.Sleep(time.Millisecond)

	user := CreateUser(pool, ctx)
	CreateOrUpdateUsersView(pool, ctx, CreateOrUpdateUsersViewStoreQuery{})

	userViews, pagination := GetUsersView(pool, ctx, GetUsersViewStoreQuery{
		Limit: 1,
		Page:  2,
	})

	require.Len(t, userViews, 1)
	require.Equal(t, int64(2), pagination.Count)
	require.False(t, pagination.HasNext)
	require.True(t, pagination.HasPrevious)
	require.Equal(t, user.Id, userViews[0].Id)
}

func BenchmarkCreateOrUpdateUsersView(b *testing.B) {
	pool := config.InitPool(nil, config.InitEnvironment())
	ctx := context.Background()
	PurgeUserViews(pool, ctx)
	PurgeUsers(pool, ctx)
	tx, _ := pool.Begin(ctx)

	for i := 1; i < 10_000; i++ {
		CreateUser(tx, ctx)
	}

	_ = tx.Commit(ctx)

	CreateOrUpdateUsersView(pool, ctx, CreateOrUpdateUsersViewStoreQuery{})
}

func TestGetUsersViewIndexUsageFilteringByIdentifier(t *testing.T) {
	pool := config.InitPool(nil, config.InitEnvironment())
	ctx := context.Background()
	SetSeqScan(pool, ctx, false)
	PurgeUserViews(pool, ctx)
	PurgeUsers(pool, ctx)
	sql := getUsersViewSQL()
	identifiers := []uuid.UUID{uuid.NewV4()}
	rows := Explain(pool, ctx, sql, functools.UUIDListToPGArray(identifiers), "{}", "{}", "{}", 1, 1)
	SetSeqScan(pool, ctx, true)
	for _, v := range rows {
		require.NotContains(t, v, "Seq Scan")
	}
}

func TestGetUsersViewIndexUsageFilteringBy(t *testing.T) {
	pool := config.InitPool(nil, config.InitEnvironment())
	ctx := context.Background()
	SetSeqScan(pool, ctx, false)
	PurgeUserViews(pool, ctx)
	PurgeUsers(pool, ctx)
	sql := getUsersViewSQL()
	identifiers := []uuid.UUID{uuid.NewV4()}
	rows := Explain(pool, ctx, sql, "{}", functools.UUIDListToPGArray(identifiers), "{}", "{}", 1, 1)
	SetSeqScan(pool, ctx, true)
	for _, v := range rows {
		require.NotContains(t, v, "Seq Scan")
	}
}
