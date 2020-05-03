package stores

import (
	"auth/models"
	"auth/repositories"
	"context"
	uuid "github.com/satori/go.uuid"
)

func (store *DatabaseStore) CreatePassword(ctx context.Context, userId uuid.UUID, value string) (int, *models.Password) {
	return repositories.CreatePassword(store.db, ctx, userId, value)
}

func (store *DatabaseStore) GetPasswords(ctx context.Context, userId uuid.UUID) []*models.Password {
	return repositories.GetPasswords(store.db, ctx, userId)
}

func (store *DatabaseStore) GetLatestPassword(ctx context.Context, userId uuid.UUID) (int, *models.Password) {
	return repositories.GetLatestPassword(store.db, ctx, userId)
}
