package stores

import (
	"auth/models"
	"auth/repositories"
	"context"
)

func (store *DatabaseStore) CreatePassword(ctx context.Context, userId int64, value string) (int, *models.Password) {
	return repositories.CreatePassword(store.Db, ctx, userId, value)
}

func (store *DatabaseStore) GetPasswords(ctx context.Context, userId int64) []*models.Password {
	return repositories.GetPasswords(store.Db, ctx, userId)
}

func (store *DatabaseStore) GetLatestPassword(ctx context.Context, userId int64) (int, *models.Password) {
	return repositories.GetLatestPassword(store.Db, ctx, userId)
}
