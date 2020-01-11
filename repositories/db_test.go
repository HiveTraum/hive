package repositories

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

func PurgeUserViews(pool *pgxpool.Pool, ctx context.Context) {
	PurgeTable(pool, ctx, "users_view")
}

func PurgeUsers(pool *pgxpool.Pool, ctx context.Context) {
	PurgeTable(pool, ctx, "users")
}

func PurgeEmails(pool *pgxpool.Pool, ctx context.Context) {
	PurgeTable(pool, ctx, "emails")
}

func PurgePhones(pool *pgxpool.Pool, ctx context.Context) {
	PurgeTable(pool, ctx, "phones")
}

func PurgeRoles(pool *pgxpool.Pool, ctx context.Context) {
	PurgeTable(pool, ctx, "roles")
}

func PurgeSessions(pool *pgxpool.Pool, ctx context.Context) {
	PurgeTable(pool, ctx, "sessions")
}

func PurgeSecrets(pool *pgxpool.Pool, ctx context.Context) {
	PurgeTable(pool, ctx, "secrets")
}

func PurgeUserRoles(pool *pgxpool.Pool, ctx context.Context) {
	PurgeTable(pool, ctx, "user_roles")
}


func PurgePasswords(pool *pgxpool.Pool, ctx context.Context) {
	PurgeTable(pool, ctx, "passwords")
}
