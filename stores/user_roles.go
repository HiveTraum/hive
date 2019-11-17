package stores

import (
	"auth/models"
	"auth/repositories"
	"context"
)

func (store *DatabaseStore) CreateUserRole(ctx context.Context, userId int64, roleId int64) (int, *models.UserRole) {
	return repositories.CreateUserRole(store.Db, ctx, userId, roleId)
}

func (store *DatabaseStore) GetUserRoles(ctx context.Context, query repositories.GetUserRoleQuery) []*models.UserRole {
	return repositories.GetUserRoles(store.Db, ctx, query)
}

func (store *DatabaseStore) DeleteUserRole(ctx context.Context, id int64) (int, *models.UserRole) {
	return repositories.DeleteUserRole(store.Db, ctx, id)
}
