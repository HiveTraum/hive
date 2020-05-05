package controllers

import (
	"hive/enums"
	"hive/models"
	"hive/repositories"
	"context"
	uuid "github.com/satori/go.uuid"
)

func (controller *Controller) CreateRole(ctx context.Context, title string) (int, *models.Role) {
	status, role := controller.store.CreateRole(ctx, title)

	if status == enums.Ok {
		controller.OnRoleChanged([]uuid.UUID{role.Id})
	}

	return status, role
}

func (controller *Controller) GetRole(ctx context.Context, id uuid.UUID) (int, *models.Role) {
	return controller.store.GetRole(ctx, id)
}

func (controller *Controller) GetRoles(ctx context.Context, query repositories.GetRolesQuery) ([]*models.Role, *models.PaginationResponse) {
	return controller.store.GetRoles(ctx, query)
}
