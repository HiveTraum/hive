package stores

import (
	"hive/models"
	"hive/repositories"
	"context"
	uuid "github.com/satori/go.uuid"
)

func (store *DatabaseStore) CreateUserRole(ctx context.Context, userId uuid.UUID, roleId uuid.UUID) (int, *models.UserRole) {
	return repositories.CreateUserRole(store.db, ctx, userId, roleId)
}

func (store *DatabaseStore) GetUserRoles(ctx context.Context, query repositories.GetUserRoleQuery) ([]*models.UserRole, *models.PaginationResponse) {
	return repositories.GetUserRoles(store.db, ctx, query)
}

func (store *DatabaseStore) DeleteUserRole(ctx context.Context, id uuid.UUID) (int, *models.UserRole) {
	return repositories.DeleteUserRole(store.db, ctx, id)
}
