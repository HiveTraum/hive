package stores

import (
	"auth/models"
	"auth/repositories"
	"context"
	uuid "github.com/satori/go.uuid"
)

func (store *DatabaseStore) CreateSession(ctx context.Context, fingerprint string, userID uuid.UUID, secretID uuid.UUID, userAgent string) (int, *models.Session) {
	return repositories.CreateSession(store.Db, ctx, fingerprint, userID, secretID, userAgent)
}

func (store *DatabaseStore) GetSession(ctx context.Context, fingerprint string, refreshToken string, userID uuid.UUID) *models.Session {
	return repositories.GetSession(store.Db, ctx, fingerprint, refreshToken, userID)
}
