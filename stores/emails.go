package stores

import (
	"auth/config"
	"auth/functools"
	"auth/models"
	"auth/repositories"
	"context"
	uuid "github.com/satori/go.uuid"
	"time"
)

func (store *DatabaseStore) CreateEmail(ctx context.Context, userId uuid.UUID, value string) (int, *models.Email) {
	return repositories.CreateEmail(store.db, ctx, userId, value)
}

func (store *DatabaseStore) GetEmail(ctx context.Context, phone string) (int, *models.Email) {
	return repositories.GetEmail(store.db, ctx, phone)
}

func (store *DatabaseStore) CreateEmailConfirmationCode(ctx context.Context, email string, code string, duration time.Duration) *models.EmailConfirmation {
	return repositories.CreateEmailConfirmationCode(ctx, store.cache, email, code, duration)
}

func (store *DatabaseStore) GetEmailConfirmationCode(ctx context.Context, email string) string {
	return repositories.GetEmailConfirmationCode(ctx, store.cache, email)
}

func (store *DatabaseStore) GetRandomCodeForEmailConfirmation() string {
	if !config.GetEnvironment().IsTestEnvironment {
		return functools.GetRandomString(6)
	} else {
		return config.GetEnvironment().TestConfirmationCode
	}
}
