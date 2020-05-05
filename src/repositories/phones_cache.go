package repositories

import (
	"hive/enums"
	"hive/models"
	"context"
	"fmt"
	"github.com/go-redis/redis/v7"
	"time"
)

func getPhoneConfirmationCodeKey(phone string) string {
	return fmt.Sprintf("%s:%s", enums.PhoneConfirmationCode, phone)
}

func CreatePhoneConfirmationCode(cache *redis.Client, ctx context.Context, phone string, code string, duration time.Duration) *models.PhoneConfirmation {

	key := getPhoneConfirmationCodeKey(phone)
	cmd := cache.WithContext(ctx).Set(key, code, duration)
	err := cmd.Err()
	if err != nil {
		return nil
	}

	created := time.Now()

	return &models.PhoneConfirmation{
		Created: created.Unix(),
		Expire:  created.Add(duration).Unix(),
		Phone:   phone,
		Code:    code,
	}

}

func GetPhoneConfirmationCode(cache *redis.Client, ctx context.Context, phone string) string {

	key := getPhoneConfirmationCodeKey(phone)
	code, err := cache.WithContext(ctx).Get(key).Result()
	if err != nil {
		return ""
	}

	return code
}
