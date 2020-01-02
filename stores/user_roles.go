package stores

import (
	"auth/models"
	"auth/repositories"
	"context"
)

func (store *DatabaseStore) CreateUserRole(ctx context.Context, userId models.UserID, roleId models.RoleID) (int, *models.UserRole) {
	return repositories.CreateUserRole(store.Db, ctx, userId, roleId)
}

func (store *DatabaseStore) GetUserRoles(ctx context.Context, query repositories.GetUserRoleQuery) []*models.UserRole {
	return repositories.GetUserRoles(store.Db, ctx, query)
}

func (store *DatabaseStore) DeleteUserRole(ctx context.Context, id models.UserRoleID) (int, *models.UserRole) {
	return repositories.DeleteUserRole(store.Db, ctx, id)
}
