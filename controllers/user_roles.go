package controllers

import (
	"auth/enums"
	"auth/infrastructure"
	"auth/models"
	"context"
	uuid "github.com/satori/go.uuid"
)

func CreateUserRole(store infrastructure.StoreInterface, esb infrastructure.ESBInterface, ctx context.Context, userId uuid.UUID, roleID uuid.UUID) (int, *models.UserRole) {
	status, userRole := store.CreateUserRole(ctx, userId, roleID)

	if status == enums.Ok {
		esb.OnUserChanged([]uuid.UUID{userRole.UserId})
	}

	return status, userRole
}

func DeleteUserRole(store infrastructure.StoreInterface, esb infrastructure.ESBInterface, ctx context.Context, id uuid.UUID) (int, *models.UserRole) {
	status, userRole := store.DeleteUserRole(ctx, id)

	if status == enums.Ok {
		esb.OnUserChanged([]uuid.UUID{userRole.UserId})
	}

	return status, userRole
}
