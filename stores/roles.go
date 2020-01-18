package stores

import (
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

func (store *DatabaseStore) GetRoles(ctx context.Context, query repositories.GetRolesQuery) []*models.Role {
	return repositories.GetRoles(store.Db, ctx, query)
}
