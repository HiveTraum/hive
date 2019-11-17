package stores

import (
	"auth/inout"
	"auth/repositories"
	"context"
)

func (store *DatabaseStore) GetUsersView(ctx context.Context, query repositories.GetUsersViewQuery) []*inout.GetUserViewResponseV1 {
	return repositories.GetUsersView(store.Db, ctx, query)
}

func (store *DatabaseStore) GetUserView(ctx context.Context, id int64) *inout.GetUserViewResponseV1 {

	userView := repositories.GetUserViewFromCache(store.Cache, ctx, id)

	if userView != nil {
		return userView
	}

	userView = repositories.GetUserView(store.Db, ctx, id)

	if userView != nil {
		store.CacheUserView(ctx, []*inout.GetUserViewResponseV1{userView})
	}

	return userView
}

func (store *DatabaseStore) CreateOrUpdateUsersView(ctx context.Context, query repositories.CreateOrUpdateUsersViewQuery) []*inout.GetUserViewResponseV1 {
	return repositories.CreateOrUpdateUsersView(store.Db, ctx, query)
}

func (store *DatabaseStore) CreateOrUpdateUsersViewByUsersID(context context.Context, id []int64) []*inout.GetUserViewResponseV1 {
	return store.CreateOrUpdateUsersView(context, repositories.CreateOrUpdateUsersViewQuery{
		GetUsersQuery: repositories.GetUsersQuery{Id: id,},
	})
}

func (store *DatabaseStore) CreateOrUpdateUsersViewByRolesID(context context.Context, id []int64) []*inout.GetUserViewResponseV1 {
	return store.CreateOrUpdateUsersView(context, repositories.CreateOrUpdateUsersViewQuery{
		GetUsersQuery: repositories.GetUsersQuery{Limit: 0, Id: nil,}, Roles: id,
	})
}

func (store *DatabaseStore) CreateOrUpdateUsersViewByUserID(context context.Context, id int64) []*inout.GetUserViewResponseV1 {
	return store.CreateOrUpdateUsersViewByUsersID(context, []int64{id})
}

func (store *DatabaseStore) CreateOrUpdateUsersViewByRoleID(context context.Context, id int64) []*inout.GetUserViewResponseV1 {
	return store.CreateOrUpdateUsersViewByRolesID(context, []int64{id})
}

// Cache

func (store *DatabaseStore) GetUserViewFromCache(ctx context.Context, id int64) *inout.GetUserViewResponseV1 {
	return repositories.GetUserViewFromCache(store.Cache, ctx, id)
}

func (store *DatabaseStore) CacheUserView(ctx context.Context, userViews []*inout.GetUserViewResponseV1) {
	repositories.CacheUserView(store.Cache, ctx, userViews)
}
