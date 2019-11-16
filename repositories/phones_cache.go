package repositories

import (
	"auth/enums"
	"auth/models"
	"context"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/opentracing/opentracing-go"
	"time"
)

func getPhoneConfirmationCodeKey(phone string) string {
	return fmt.Sprintf("%s:%s", enums.PhoneConfirmationCode, phone)
}

func CreatePhoneConfirmationCode(cache *redis.Client, ctx context.Context, phone string, code string, duration time.Duration) *models.PhoneConfirmation {

	span, ctx := opentracing.StartSpanFromContext(ctx, "Create phone confirmation code")

	key := getPhoneConfirmationCodeKey(phone)
	cache.Set(key, code, duration)
	created := time.Now()

	span.Finish()

	return &models.PhoneConfirmation{
		Created: created.Unix(),
		Expire:  created.Add(duration).Unix(),
		Phone:   phone,
		Code:    code,
	}

}

func GetPhoneConfirmationCode(cache *redis.Client, ctx context.Context, phone string) (string, error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "Get phone confirmation code")

	key := getPhoneConfirmationCodeKey(phone)
	code, err := cache.Get(key).Result()
	if err != nil {
		return "", err
	}

	span.Finish()

	return code, nil
}
