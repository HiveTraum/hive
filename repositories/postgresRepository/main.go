package postgresRepository

import (
	"auth/config"
	"auth/models"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	uuid "github.com/satori/go.uuid"
)

type IPostgresRepository interface {

	// Secrets

	CreateSecret(ctx context.Context) *models.Secret
	GetSecret(ctx context.Context, id uuid.UUID) *models.Secret

	// Sessions

	CreateSession(ctx context.Context, userID uuid.UUID, secretID uuid.UUID, fingerprint string, userAgent string) *models.Session
	DeleteSession(ctx context.Context, id uuid.UUID) *models.Session
}

type PostgresRepository struct {
	pool        *pgxpool.Pool
	environment *config.Environment
}

func InitPostgresRepository(pool *pgxpool.Pool, environment *config.Environment) *PostgresRepository {
	return &PostgresRepository{
		pool:        pool,
		environment: environment,
	}
}
