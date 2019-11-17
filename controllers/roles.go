package controllers

import (
	"auth/enums"
	"auth/infrastructure"
	"auth/models"
	"context"
)

func CreateRole(store infrastructure.StoreInterface, esb infrastructure.ESBInterface, ctx context.Context, title string) (int, *models.Role) {
	status, role := store.CreateRole(ctx, title)

	if status == enums.Ok {
		esb.OnRoleChanged([]int64{role.Id})
	}

	return status, role
}
