package controllers

import (
	"auth/enums"
	"auth/infrastructure"
	"auth/models"
	"context"
)

func CreateUserRole(store infrastructure.StoreInterface, esb infrastructure.ESBInterface, ctx context.Context, userId models.UserID, roleId models.RoleID) (int, *models.UserRole) {
	status, userRole := store.CreateUserRole(ctx, userId, roleId)

	if status == enums.Ok {
		esb.OnUserChanged([]models.UserID{userRole.UserId})
	}

	return status, userRole
}

func DeleteUserRole(store infrastructure.StoreInterface, esb infrastructure.ESBInterface, ctx context.Context, id models.UserRoleID) (int, *models.UserRole) {
	status, userRole := store.DeleteUserRole(ctx, id)

	if status == enums.Ok {
		esb.OnUserChanged([]models.UserID{userRole.UserId})
	}

	return status, userRole
}
