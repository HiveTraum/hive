package repositories

import (
	"auth/enums"
	"auth/models"
	"context"
	"errors"
	"github.com/getsentry/sentry-go"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	uuid "github.com/satori/go.uuid"
	"strings"
)

func createPhoneSQL() string {
	return `INSERT INTO phones (id, user_id, value, country_code) 
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (value) 
			    DO UPDATE SET created=DEFAULT,
			                  user_id=excluded.user_id
			RETURNING id, created, user_id, value, country_code;`
}

func getPhoneSQL() string {
	return `SELECT id, created, user_id, value, country_code FROM phones WHERE value = $1;`
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

	err := row.Scan(&phone.Id, &phone.Created, &phone.UserId, &phone.Value, &phone.CountryCode)
	if err != nil {
		return unwrapPhoneScanError(err), nil
	}

	return enums.Ok, phone
}

func CreatePhone(db DB, ctx context.Context, userId uuid.UUID, value string, countryCode string) (int, *models.Phone) {
	sql := createPhoneSQL()
	row := db.QueryRow(ctx, sql, uuid.NewV4(), userId, value, countryCode)
	return scanPhone(row)
}

func GetPhone(db DB, ctx context.Context, phone string) (int, *models.Phone) {
	sql := getPhoneSQL()
	row := db.QueryRow(ctx, sql, phone)
	return scanPhone(row)
}
