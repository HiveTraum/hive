package repositories

import (
	"auth/functools"
	"auth/models"
	"context"
	"github.com/getsentry/sentry-go"
	"github.com/jackc/pgx/v4"
	uuid "github.com/satori/go.uuid"
)

func createSecretSQL() string {
	return "INSERT INTO secrets (id, created, value) VALUES ($1, default, default) RETURNING id, created, value;"
}

func getSecretsSQL() string {
	return `
		SELECT id, created, value
		FROM secrets
		WHERE (array_length($1::uuid[], 1) IS NULL OR id = ANY ($1::uuid[]))
		LIMIT $2;
		`
}

func scanSecret(row pgx.Row) *models.Secret {
	secret := &models.Secret{}

	err := row.Scan(&secret.Id, &secret.Created, &secret.Value)
	if err != nil {
		sentry.CaptureException(err)
		return nil
	}

	return secret
}

func CreateSecret(db DB, ctx context.Context) *models.Secret {
	sql := createSecretSQL()
	row := db.QueryRow(ctx, sql, uuid.NewV4())
	return scanSecret(row)
}

func GetSecretFromDB(db DB, ctx context.Context, id uuid.UUID) *models.Secret {
	sql := getSecretsSQL()
	row := db.QueryRow(ctx, sql, functools.StringsToPGArray([]string{id.String()}), 1)
	return scanSecret(row)
}
