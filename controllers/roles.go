package controllers

import (
	"auth/infrastructure"
	"auth/models"
	"context"
)

func CreateRole(store infrastructure.StoreInterface, esb infrastructure.ESBInterface, ctx context.Context, title string) *models.Role {
	role := store.CreateRole(ctx, title)
	esb.OnRoleChanged([]int64{role.Id})
	return role
}
