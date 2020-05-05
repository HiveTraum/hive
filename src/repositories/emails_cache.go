package repositories

import (
	"hive/enums"
	"hive/models"
	"context"
	"fmt"
	"github.com/go-redis/redis/v7"
	"time"
)

func getEmailConfirmationCodeKey(email string) string {
	return fmt.Sprintf("%s:%s", enums.EmailConfirmationCode, email)
}

func CreateEmailConfirmationCode(ctx context.Context, cache *redis.Client, email string, code string, duration time.Duration) *models.EmailConfirmation {
	key := getEmailConfirmationCodeKey(email)
	cmd := cache.WithContext(ctx).Set(key, code, duration)
	if err := cmd.Err(); err != nil {
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

func GetEmailConfirmationCode(ctx context.Context, cache *redis.Client, email string) string {
	key := getEmailConfirmationCodeKey(email)
	code, err := cache.WithContext(ctx).Get(key).Result()

	if err != nil {
		return ""
	}

	return code
}
