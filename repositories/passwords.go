package repositories

import (
	"auth/enums"
	"auth/models"
	"context"
	"errors"
	"github.com/getsentry/sentry-go"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"strings"
)

func createPasswordSQL() string {
	return `INSERT INTO passwords (user_id, value) 
			VALUES ($1, $2) 
			RETURNING id, created, user_id, value;`
}

func getPasswordSQL() string {
	return `SELECT id, created, user_id, value 
			FROM passwords 
			WHERE user_id = $1 
			LIMIT $2;`
}

func getLatestPasswordSQL() string {
	return `SELECT id, created, user_id, value
			FROM passwords
			WHERE user_id = $1
			ORDER BY id DESC
			LIMIT 1;`
}

func unwrapPasswordScanErrors(err error) int {
	var e *pgconn.PgError
	if errors.As(err, &e) && strings.Contains(e.Detail, "is not present in table \"users\"") {
		return enums.UserNotFound
	}

	sentry.CaptureException(err)
	return enums.NotOk
}

func scanPassword(row pgx.Row) (int, *models.Password) {
	password := &models.Password{}

	err := row.Scan(&password.Id, &password.Created, &password.UserId, &password.Value)
	if err != nil {
		return unwrapPasswordScanErrors(err), nil
	}

	return enums.Ok, password
}

func scanPasswords(rows pgx.Rows, limit int) []*models.Password {
	passwords := make([]*models.Password, limit)

	var i int32

	for rows.Next() {
		_, password := scanPassword(rows)
		passwords[i] = password
		i++
	}

	rows.Close()

	return passwords[0:i]
}

func CreatePassword(db DB, ctx context.Context, userId models.UserID, value string) (int, *models.Password) {
	sql := createPasswordSQL()
	row := db.QueryRow(ctx, sql, userId, value)
	return scanPassword(row)
}

func GetPasswords(db DB, ctx context.Context, userId models.UserID) []*models.Password {
	sql := getPasswordSQL()
	limit := 10
	rows, err := db.Query(ctx, sql, userId, limit)
	if err != nil {
		sentry.CaptureException(err)
		return nil
	}

	return scanPasswords(rows, limit)
}

func GetLatestPassword(db DB, ctx context.Context, userId models.UserID) (int, *models.Password) {
	sql := getLatestPasswordSQL()
	row := db.QueryRow(ctx, sql, userId)
	return scanPassword(row)
}
