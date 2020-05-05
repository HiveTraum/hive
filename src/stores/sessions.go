package stores

import (
	"hive/models"
	"context"
	uuid "github.com/satori/go.uuid"
)

func (store *DatabaseStore) CreateSession(ctx context.Context, userID uuid.UUID, secretID uuid.UUID, fingerprint string, userAgent string) *models.Session {
	return store.postgresRepository.CreateSession(ctx, userID, secretID, fingerprint, userAgent)
}

func (store *DatabaseStore) DeleteSession(ctx context.Context, id uuid.UUID) *models.Session {
	return store.postgresRepository.DeleteSession(ctx, id)
}
