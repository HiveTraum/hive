package repositories

import (
	"auth/enums"
	"auth/models"
	"context"
	"github.com/getsentry/sentry-go"
	"github.com/jackc/pgx/v4"
)

func createSessionSQL() string {
	return `INSERT INTO sessions (fingerprint, user_id, secret_id, user_agent) 
			VALUES ($1, $2, $3, $4)
			RETURNING refresh_token, fingerprint, user_id, secret_id, created, user_agent;`
}

func getSessionsSQL() string {
	return `
			SELECT refresh_token, fingerprint, user_id, secret_id, created, user_agent
			FROM sessions
			WHERE refresh_token = $1 AND fingerprint = $2 AND user_id = $3
			LIMIT $4;
			`
}

func scanSession(row pgx.Row) *models.Session {
	session := &models.Session{}

	err := row.Scan(&session.RefreshToken, &session.Fingerprint, &session.UserID, &session.SecretID, &session.Created, &session.UserAgent)
	if err != nil {
		sentry.CaptureException(err)
		return nil
	}

	return session
}

func CreateSession(db DB, ctx context.Context, fingerprint string, userID models.UserID, secretID models.SecretID, userAgent string) (int, *models.Session) {
	sql := createSessionSQL()
	row := db.QueryRow(ctx, sql, fingerprint, userID, secretID, userAgent)
	return enums.Ok, scanSession(row)
}

func GetSession(db DB, ctx context.Context, fingerprint string, refreshToken string, userID models.UserID) *models.Session {
	sql := getSessionsSQL()
	row := db.QueryRow(ctx, sql, refreshToken, fingerprint, userID, 1)
	return scanSession(row)
}
