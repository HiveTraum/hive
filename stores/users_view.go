package stores

import (
	"auth/models"
	"auth/repositories"
	"context"
	uuid "github.com/satori/go.uuid"
)

func (store *DatabaseStore) GetUsersView(ctx context.Context, query repositories.GetUsersViewStoreQuery) ([]*models.UserView, *models.PaginationResponse) {
	return repositories.GetUsersView(store.db, ctx, query)
}

func (store *DatabaseStore) GetUserView(ctx context.Context, id uuid.UUID) *models.UserView {

	userView := repositories.GetUserViewFromCache(store.cache, ctx, id)

	if userView != nil {
		return userView
	}

	userView = repositories.GetUserView(store.db, ctx, id)

	if userView != nil {
		store.CacheUserView(ctx, []*models.UserView{userView})
	}

	return userView
}

func (store *DatabaseStore) CreateOrUpdateUsersView(ctx context.Context, query repositories.CreateOrUpdateUsersViewStoreQuery) []*models.UserView {
	return repositories.CreateOrUpdateUsersView(store.db, ctx, query)
}

func (store *DatabaseStore) CreateOrUpdateUsersViewByUsersID(context context.Context, id []uuid.UUID) []*models.UserView {
	return store.CreateOrUpdateUsersView(context, repositories.CreateOrUpdateUsersViewStoreQuery{Id: id,})
}

func (store *DatabaseStore) CreateOrUpdateUsersViewByRolesID(context context.Context, id []uuid.UUID) []*models.UserView {
	return store.CreateOrUpdateUsersView(context, repositories.CreateOrUpdateUsersViewStoreQuery{
		Limit: 0, Id: nil, Roles: id,
	})
}

func (store *DatabaseStore) CreateOrUpdateUsersViewByUserID(context context.Context, id uuid.UUID) []*models.UserView {
	return store.CreateOrUpdateUsersViewByUsersID(context, []uuid.UUID{id})
}

func (store *DatabaseStore) CreateOrUpdateUsersViewByRoleID(context context.Context, id uuid.UUID) []*models.UserView {
	return store.CreateOrUpdateUsersViewByRolesID(context, []uuid.UUID{id})
}

// cache

func (store *DatabaseStore) GetUserViewFromCache(ctx context.Context, id uuid.UUID) *models.UserView {
	return repositories.GetUserViewFromCache(store.cache, ctx, id)
}

func (store *DatabaseStore) CacheUserView(ctx context.Context, userViews []*models.UserView) {

	repositories.CacheUserView(store.cache, ctx, userViews)
}
