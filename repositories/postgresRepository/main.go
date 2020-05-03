package postgresRepository

import (
	"auth/models"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	uuid "github.com/satori/go.uuid"
)

type IPostgresRepository interface {

	// Secrets

	CreateSecret(ctx context.Context) *models.Secret
	GetSecret(ctx context.Context, id uuid.UUID) *models.Secret
}

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func InitPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		pool: pool,
	}
}
