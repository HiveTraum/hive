package postgresRepository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
)

func PurgeTable(pool *pgxpool.Pool, ctx context.Context, table string) {
	_, err := pool.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE ", table))
	if err != nil {
		panic(err)
	}
}

func PurgeSecrets(pool *pgxpool.Pool, ctx context.Context) {
	PurgeTable(pool, ctx, "secrets")
}
