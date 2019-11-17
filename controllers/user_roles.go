package controllers

import (
	"auth/enums"
	"auth/infrastructure"
	"auth/models"
	"context"
)

func CreateUserRole(store infrastructure.StoreInterface, esb infrastructure.ESBInterface, ctx context.Context, userId int64, roleId int64) (int, *models.UserRole) {
	status, userRole := store.CreateUserRole(ctx, userId, roleId)

	if status == enums.Ok {
		esb.OnUserChanged([]int64{userRole.UserId})
	}

	return status, userRole
}

func DeleteUserRole(store infrastructure.StoreInterface, esb infrastructure.ESBInterface, ctx context.Context, id int64) (int, *models.UserRole) {
	status, userRole := store.DeleteUserRole(ctx, id)

	if status == enums.Ok {
		esb.OnUserChanged([]int64{userRole.UserId})
	}

	return status, userRole
}
