package repositories

import (
	"auth/enums"
	"auth/models"
	"context"
	"errors"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/go-redis/redis/v7"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"strings"
	"time"
)

func createPhoneSQL() string {
	return `INSERT INTO phones (user_id, value) 
			VALUES ($1, $2)
			ON CONFLICT (value) 
			    DO UPDATE SET created=DEFAULT,
			                  user_id=excluded.user_id
			RETURNING id, created, user_id, value;`
}

func getPhoneSQL() string {
	return `SELECT id, created, user_id, value FROM phones WHERE value = $1;`
}

func unwrapPhoneScanError(err error) int {
	var e *pgconn.PgError
	if errors.As(err, &e) && strings.Contains(e.Detail, "is not present in table \"users\"") {
		return enums.UserNotFound
	} else if strings.Contains(err.Error(), "no rows") {
		return enums.Ok
	}

	sentry.CaptureException(err)
	return enums.NotOk
}

func scanPhone(row pgx.Row) (int, *models.Phone) {
	phone := &models.Phone{}

	err := row.Scan(&phone.Id, &phone.Created, &phone.UserId, &phone.Value)
	if err != nil {
		return unwrapPhoneScanError(err), nil
	}

	return enums.Ok, phone
}

func getPhoneConfirmationCodeKey(phone string) string {
	return fmt.Sprintf("%s:%s", enums.PhoneConfirmationCode, phone)
}

func CreatePhone(db DB, ctx context.Context, userId int64, value string) (int, *models.Phone) {
	sql := createPhoneSQL()
	row := db.QueryRow(ctx, sql, userId, value)
	return scanPhone(row)
}

func GetPhone(db DB, ctx context.Context, phone string) (int, *models.Phone) {
	sql := getPhoneSQL()
	row := db.QueryRow(ctx, sql, phone)
	return scanPhone(row)
}

func CreatePhoneConfirmationCode(cache *redis.Client, phone string, code string, duration time.Duration) *models.PhoneConfirmation {
	key := getPhoneConfirmationCodeKey(phone)
	cache.Set(key, code, duration)
	created := time.Now()
	return &models.PhoneConfirmation{
		Created: created.Unix(),
		Expire:  created.Add(duration).Unix(),
		Phone:   phone,
		Code:    code,
	}
}

func GetPhoneConfirmationCode(cache *redis.Client, phone string) (string, error) {
	key := getPhoneConfirmationCodeKey(phone)
	code, err := cache.Get(key).Result()
	if err != nil {
		return "", err
	}

	return code, nil
}
