package stores

import (
	"auth/models"
	"auth/repositories"
	"context"
	"github.com/getsentry/sentry-go"
	"time"
)

func (store *DatabaseStore) CreatePhone(ctx context.Context, userId int64, value string) (int, *models.Phone) {
	return repositories.CreatePhone(store.Db, ctx, userId, value)
}

func (store *DatabaseStore) GetPhone(ctx context.Context, phone string) (int, *models.Phone) {
	return repositories.GetPhone(store.Db, ctx, phone)
}

func (store *DatabaseStore) CreatePhoneConfirmationCode(ctx context.Context, phone string, code string, duration time.Duration) *models.PhoneConfirmation {
	return repositories.CreatePhoneConfirmationCode(store.Cache, ctx, phone, code, duration)
}

func (store *DatabaseStore) GetPhoneConfirmationCode(ctx context.Context, phone string) string {
	code, err := repositories.GetPhoneConfirmationCode(store.Cache, ctx, phone)

	if err != nil {
		sentry.CaptureException(err)
	}

	return code
}
