package stores

import (
	"auth/models"
	"auth/repositories"
	"context"
	"github.com/getsentry/sentry-go"
	"time"
)

func (store *DatabaseStore) CreateEmail(ctx context.Context, userId models.UserID, value string) (int, *models.Email) {
	return repositories.CreateEmail(store.Db, ctx, userId, value)
}

func (store *DatabaseStore) GetEmail(ctx context.Context, phone string) (int, *models.Email) {
	return repositories.GetEmail(store.Db, ctx, phone)
}

func (store *DatabaseStore) CreateEmailConfirmationCode(email string, code string, duration time.Duration) *models.EmailConfirmation {
	return repositories.CreateEmailConfirmationCode(store.Cache, email, code, duration)
}

func (store *DatabaseStore) GetEmailConfirmationCode(email string) string {
	code, err := repositories.GetEmailConfirmationCode(store.Cache, email)

	if err != nil {
		sentry.CaptureException(err)
	}

	return code
}
