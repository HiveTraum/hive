package controllers

import (
	"auth/infrastructure"
	"auth/models"
	"context"
	uuid "github.com/satori/go.uuid"
)

func CreateOrUpdateUsersView(store infrastructure.StoreInterface, esb infrastructure.ESBInterface, ctx context.Context, id []uuid.UUID) []*models.UserView {
	usersView := store.CreateOrUpdateUsersViewByUsersID(ctx, id)
	esb.OnUsersViewChanged(usersView)
	return usersView
}

func CreateOrUpdateUsersViewByRoles(store infrastructure.StoreInterface, esb infrastructure.ESBInterface, ctx context.Context, rolesIds []uuid.UUID) []*models.UserView {
	usersView := store.CreateOrUpdateUsersViewByRolesID(ctx, rolesIds)
	esb.OnUsersViewChanged(usersView)
	return usersView
}
