package controllers

import (
	"hive/models"
	"hive/repositories"
	"context"
	uuid "github.com/satori/go.uuid"
)

func (controller *Controller) CreateOrUpdateUsersView(ctx context.Context, id []uuid.UUID) []*models.UserView {
	usersView := controller.store.CreateOrUpdateUsersViewByUsersID(ctx, id)
	controller.OnUsersViewChanged(usersView)
	return usersView
}

func (controller *Controller) CreateOrUpdateUsersViewByRoles(ctx context.Context, rolesIds []uuid.UUID) []*models.UserView {
	usersView := controller.store.CreateOrUpdateUsersViewByRolesID(ctx, rolesIds)
	controller.OnUsersViewChanged(usersView)
	return usersView
}

func (controller *Controller) GetUserView(ctx context.Context, id uuid.UUID) *models.UserView {
	return controller.store.GetUserView(ctx, id)
}

func (controller *Controller) GetUserViews(ctx context.Context, query repositories.GetUsersViewStoreQuery) ([]*models.UserView, *models.PaginationResponse) {
	return controller.store.GetUsersView(ctx, query)
}
