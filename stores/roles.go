package stores

import (
	"auth/config"
	"auth/enums"
	"auth/models"
	"auth/repositories"
	"context"
	uuid "github.com/satori/go.uuid"
)

func (store *DatabaseStore) CreateRole(ctx context.Context, title string) (int, *models.Role) {
	return repositories.CreateRole(store.Db, ctx, title)
}

func (store *DatabaseStore) GetRole(ctx context.Context, id uuid.UUID) (int, *models.Role) {
	return repositories.GetRole(store.Db, ctx, id)
}

func (store *DatabaseStore) GetRoles(ctx context.Context, query repositories.GetRolesQuery) ([]*models.Role, *models.PaginationResponse) {
	return repositories.GetRoles(store.Db, ctx, query)
}

func (store *DatabaseStore) GetRoleByTitle(ctx context.Context, title string) (int, *models.Role) {
	roles, _ := store.GetRoles(ctx, repositories.GetRolesQuery{
		Pagination: &models.PaginationRequest{
			Page:  1,
			Limit: 1,
		},
		Titles: []string{title},
	})

	if len(roles) > 0 {
		return enums.Ok, roles[0]
	} else {
		return enums.Ok, nil
	}
}

func (store *DatabaseStore) GetAdminRole(ctx context.Context) (int, *models.Role) {
	return store.GetRoleByTitle(ctx, config.GetEnvironment().AdminRole)
}
