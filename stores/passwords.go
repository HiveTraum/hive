package stores

import (
	"auth/models"
	"auth/repositories"
	"context"
)

func (store *DatabaseStore) CreatePassword(ctx context.Context, userId models.UserID, value string) (int, *models.Password) {
	return repositories.CreatePassword(store.Db, ctx, userId, value)
}

func (store *DatabaseStore) GetPasswords(ctx context.Context, userId models.UserID) []*models.Password {
	return repositories.GetPasswords(store.Db, ctx, userId)
}

func (store *DatabaseStore) GetLatestPassword(ctx context.Context, userId models.UserID) (int, *models.Password) {
	return repositories.GetLatestPassword(store.Db, ctx, userId)
}
