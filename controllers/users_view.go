package controllers

import (
	"auth/infrastructure"
	"auth/inout"
	"context"
)

func CreateOrUpdateUsersView(store infrastructure.StoreInterface, esb infrastructure.ESBInterface, ctx context.Context, id []int64) []*inout.GetUserViewResponseV1 {
	usersView := store.CreateOrUpdateUsersViewByUsersID(ctx, id)
	esb.OnUsersViewChanged(usersView)
	return usersView
}

func CreateOrUpdateUsersViewByRoles(store infrastructure.StoreInterface, esb infrastructure.ESBInterface, ctx context.Context, rolesIds []int64) []*inout.GetUserViewResponseV1 {
	usersView := store.CreateOrUpdateUsersViewByRolesID(ctx, rolesIds)
	esb.OnUsersViewChanged(usersView)
	return usersView
}

func GetUserView(store infrastructure.StoreInterface, esb infrastructure.ESBInterface, ctx context.Context, id int64) *inout.GetUserViewResponseV1 {
	userView := store.GetUserViewFromCache(id)

	if userView != nil {
		return userView
	}

	userView = store.GetUserView(ctx, id)

	if userView != nil {
		store.CacheUserView([]*inout.GetUserViewResponseV1{userView})
	}

	return userView
}
