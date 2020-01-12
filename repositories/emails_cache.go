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
	cmd := cache.Set(key, code, duration)
	err := cmd.Err()

	if err != nil {
		sentry.CaptureException(err)
		return nil
	}

	created := time.Now()
	return &models.EmailConfirmation{
		Created: created.Unix(),
		Expire:  created.Add(duration).Unix(),
		Email:   email,
		Code:    code,
	}
}

func GetEmailConfirmationCode(cache *redis.Client, email string) string {
	key := getEmailConfirmationCodeKey(email)
	code, err := cache.Get(key).Result()

	if err != nil {
		if err != redis.Nil {
			sentry.CaptureException(err)
		}

		return ""
	}

	return code
}
