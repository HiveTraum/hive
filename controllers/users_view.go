package controllers

import (
	"auth/infrastructure"
	"auth/inout"
	"auth/models"
	"context"
)

func CreateOrUpdateUsersView(store infrastructure.StoreInterface, esb infrastructure.ESBInterface, ctx context.Context, id []models.UserID) []*inout.GetUserViewResponseV1 {
	usersView := store.CreateOrUpdateUsersViewByUsersID(ctx, id)
	esb.OnUsersViewChanged(usersView)
	return usersView
}

func CreateOrUpdateUsersViewByRoles(store infrastructure.StoreInterface, esb infrastructure.ESBInterface, ctx context.Context, rolesIds []models.RoleID) []*inout.GetUserViewResponseV1 {
	usersView := store.CreateOrUpdateUsersViewByRolesID(ctx, rolesIds)
	esb.OnUsersViewChanged(usersView)
	return usersView
}
