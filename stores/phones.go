package stores

import (
	"auth/config"
	"auth/models"
	"auth/repositories"
	"context"
	uuid "github.com/satori/go.uuid"
	"math/rand"
	"strconv"
	"time"
)

func (store *DatabaseStore) CreatePhone(ctx context.Context, userId uuid.UUID, value string) (int, *models.Phone) {
	return repositories.CreatePhone(store.Db, ctx, userId, value)
}

func (store *DatabaseStore) GetPhone(ctx context.Context, phone string) (int, *models.Phone) {
	return repositories.GetPhone(store.Db, ctx, phone)
}

func (store *DatabaseStore) CreatePhoneConfirmationCode(ctx context.Context, phone string, code string, duration time.Duration) *models.PhoneConfirmation {
	return repositories.CreatePhoneConfirmationCode(store.Cache, ctx, phone, code, duration)
}

func (store *DatabaseStore) GetPhoneConfirmationCode(ctx context.Context, phone string) string {
	return repositories.GetPhoneConfirmationCode(store.Cache, ctx, phone)
}

func (store *DatabaseStore) GetRandomCodeForPhoneConfirmation() string {

	if !config.GetEnvironment().IsTestEnvironment {
		rand.Seed(time.Now().UnixNano())
		min := 100000
		max := 999999
		return strconv.Itoa(rand.Intn(max-min+1) + min)
	} else {
		return config.GetEnvironment().TestConfirmationCode
	}
}
