package repositories

import (
	"auth/models"
	"auth/modelsFunctools"
	"context"
	"github.com/getsentry/sentry-go"
	"github.com/jackc/pgx/v4"
)

func createSecretSQL() string {
	return "INSERT INTO secrets DEFAULT VALUES RETURNING id, created, value;"
}

func getSecretsSQL() string {
	return `
		SELECT id, created, value
		FROM secrets
		WHERE (array_length($1::integer[], 1) IS NULL OR id = ANY ($1::bigint[]))
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
	row := db.QueryRow(ctx, sql)
	return scanSecret(row)
}

func GetSecretFromDB(db DB, ctx context.Context, id models.SecretID) *models.Secret {
	sql := getSecretsSQL()
	row := db.QueryRow(ctx, sql, modelsFunctools.SecretIDListToPGArray([]models.SecretID{id}), 1)
	return scanSecret(row)
}
