package controllers

import (
	"hive/enums"
	"hive/models"
	"hive/repositories"
	"context"
	uuid "github.com/satori/go.uuid"
)

func (controller *Controller) CreateUserRole(ctx context.Context, userId uuid.UUID, roleID uuid.UUID) (int, *models.UserRole) {
	status, userRole := controller.store.CreateUserRole(ctx, userId, roleID)

	if status == enums.Ok {
		controller.OnUserChanged([]uuid.UUID{userRole.UserId})
	}

	return status, userRole
}

func (controller *Controller) GetUserRoles(ctx context.Context, query repositories.GetUserRoleQuery) ([]*models.UserRole, *models.PaginationResponse) {
	return controller.store.GetUserRoles(ctx, query)
}

func (controller *Controller) DeleteUserRole(ctx context.Context, id uuid.UUID) (int, *models.UserRole) {
	status, userRole := controller.store.DeleteUserRole(ctx, id)

	if status == enums.Ok {
		controller.OnUserChanged([]uuid.UUID{userRole.UserId})
	}

	return status, userRole
}
