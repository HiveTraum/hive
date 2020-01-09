package stores

import (
	"auth/models"
	"auth/repositories"
	"context"
)

func (store *DatabaseStore) CreateSession(ctx context.Context, fingerprint string, userID models.UserID, secretID models.SecretID, userAgent string) (int, *models.Session) {
	return repositories.CreateSession(store.Db, ctx, fingerprint, userID, secretID, userAgent)
}

func (store *DatabaseStore) GetSession(ctx context.Context, fingerprint string, refreshToken string) *models.Session {
	return repositories.GetSession(store.Db, ctx, fingerprint, refreshToken)
}
