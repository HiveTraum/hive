package stores

import (
	"auth/models"
	"auth/repositories"
	"context"
	"time"
)

func (store *DatabaseStore) CreateEmail(ctx context.Context, userId models.UserID, value string) (int, *models.Email) {
	return repositories.CreateEmail(store.Db, ctx, userId, value)
}

func (store *DatabaseStore) GetEmail(ctx context.Context, phone string) (int, *models.Email) {
	return repositories.GetEmail(store.Db, ctx, phone)
}

func (store *DatabaseStore) CreateEmailConfirmationCode(_ context.Context, email string, code string, duration time.Duration) *models.EmailConfirmation {
	return repositories.CreateEmailConfirmationCode(store.Cache, email, code, duration)
}

func (store *DatabaseStore) GetEmailConfirmationCode(_ context.Context, email string) string {
	return repositories.GetEmailConfirmationCode(store.Cache, email)
}
