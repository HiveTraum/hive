package repositories

import (
	"auth/enums"
	"auth/models"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/go-redis/redis/v7"
	"time"
)

func getEmailConfirmationCodeKey(email string) string {
	return fmt.Sprintf("%s:%s", enums.EmailConfirmationCode, email)
}

func CreateEmailConfirmationCode(cache *redis.Client, email string, code string, duration time.Duration) *models.EmailConfirmation {
	key := getEmailConfirmationCodeKey(email)
	cache.Set(key, code, duration)
	created := time.Now()
	return &models.EmailConfirmation{
		Created: created.Unix(),
		Expire:  created.Add(duration).Unix(),
		Email:   email,
		Code:    code,
	}
}

func GetEmailConfirmationCode(cache *redis.Client, email string) (string, error) {
	key := getEmailConfirmationCodeKey(email)
	code, err := cache.Get(key).Result()

	if err != nil {
		sentry.CaptureException(err)
		return "", err
	}

	return code, nil
}