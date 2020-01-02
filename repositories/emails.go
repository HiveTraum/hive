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

func createEmailSQL() string {
	return `INSERT INTO emails (user_id, value) 
			VALUES ($1, $2)
			ON CONFLICT (value)
			    DO UPDATE SET created=DEFAULT,
			                  user_id=excluded.user_id
			RETURNING id, created, user_id, value;`
}

func getEmailSQL() string {
	return `SELECT id, created, user_id, value FROM emails WHERE value = $1;`
}

func unwrapEmailScanError(err error) int {
	var e *pgconn.PgError
	if errors.As(err, &e) && strings.Contains(e.Detail, "is not present in table \"users\"") {
		return enums.UserNotFound
	} else if strings.Contains(err.Error(), "no rows") {
		return enums.Ok
	}

	sentry.CaptureException(err)
	return enums.NotOk
}

func scanEmail(row pgx.Row) (int, *models.Email) {
	email := &models.Email{}

	err := row.Scan(&email.Id, &email.Created, &email.UserId, &email.Value)
	if err != nil {
		return unwrapEmailScanError(err), nil
	}

	return enums.Ok, email
}

func getEmailConfirmationCodeKey(email string) string {
	return fmt.Sprintf("%s:%s", enums.EmailConfirmationCode, email)
}

func CreateEmail(db DB, ctx context.Context, userId models.UserID, value string) (int, *models.Email) {
	sql := createEmailSQL()
	row := db.QueryRow(ctx, sql, userId, value)
	return scanEmail(row)
}

func GetEmail(db DB, ctx context.Context, phone string) (int, *models.Email) {
	sql := getEmailSQL()
	row := db.QueryRow(ctx, sql, phone)
	return scanEmail(row)
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
