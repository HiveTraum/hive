package postgresRepository

import (
	"hive/functools"
	"hive/models"
	"context"
	"github.com/getsentry/sentry-go"
	"github.com/jackc/pgx/v4"
	uuid "github.com/satori/go.uuid"
	"time"
)

func createSessionSQL() string {
	return `INSERT INTO sessions (id, user_id, secret_id, fingerprint, user_agent, created, expires) 
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id, user_id, secret_id, fingerprint, user_agent, created, expires;`
}

func deleteSessionSQL() string {
	return `
			DELETE FROM sessions 
			WHERE id = $1::uuid
			RETURNING id, user_id, secret_id, fingerprint, user_agent, created, expires;
			`
}

func getSessionsSQL() string {
	return `
			SELECT id, user_id, secret_id, fingerprint, user_agent, created, expires
			FROM sessions
			WHERE id = $1::uuid;
			`
}

func scanSession(row pgx.Row) *models.Session {
	session := &models.Session{}

	err := row.Scan(
		&session.Id,
		&session.UserID,
		&session.SecretID,
		&session.Fingerprint,
		&session.UserAgent,
		&session.Created,
		&session.Expires)
	if err != nil {
		sentry.CaptureException(err)
		return nil
	}

	return session
}

func scanSessions(rows pgx.Rows, limit int) []*models.Session {
	sessions := make([]*models.Session, limit)

	var i int32

	for rows.Next() {
		session := scanSession(rows)
		sessions[i] = session
		i++
	}

	return sessions[0:i]
}

type GetSessionsQuery struct {
	Pagination      *models.PaginationRequest
	Identifiers     []uuid.UUID
	Fingerprints    []string
	RefreshTokens   []uuid.UUID
	UserIdentifiers []uuid.UUID
	UserAgents      []string
	Now             int64
}

func (repository *PostgresRepository) CreateSession(ctx context.Context, userID, secretID uuid.UUID, fingerprint, userAgent string) *models.Session {
	sql := createSessionSQL()
	created := time.Now()
	expires := time.Now().Add(time.Minute * time.Duration(repository.environment.AccessTokenLifetime))
	row := repository.pool.QueryRow(ctx, sql, uuid.NewV4(), userID, secretID, fingerprint, userAgent, created.Unix(), expires.Unix())
	return scanSession(row)
}

func (repository *PostgresRepository) GetSessions(ctx context.Context, query *GetSessionsQuery) []*models.Session {
	sql := getSessionsSQL()
	row, _ := repository.pool.Query(ctx, sql,
		functools.UUIDListToPGArray(query.Identifiers),
		functools.UUIDListToPGArray(query.UserIdentifiers),
		functools.UUIDListToPGArray(query.RefreshTokens),
		functools.StringsToPGArray(query.Fingerprints),
		functools.StringsToPGArray(query.UserAgents),
		query.Now,
		query.Pagination.Limit)
	return scanSessions(row, query.Pagination.Limit)
}

func (repository *PostgresRepository) GetSession(ctx context.Context, id uuid.UUID) *models.Session {
	sql := getSessionsSQL()
	row := repository.pool.QueryRow(ctx, sql, id)
	return scanSession(row)
}

func (repository *PostgresRepository) DeleteSession(ctx context.Context, id uuid.UUID) *models.Session {
	sql := deleteSessionSQL()
	row := repository.pool.QueryRow(ctx, sql, id)
	return scanSession(row)
}
