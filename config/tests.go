package config

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
)

func PurgeTable(pool *pgxpool.Pool, ctx context.Context, table string) {
	_, err := pool.Exec(ctx, fmt.Sprintf("TRUNCATE %s CASCADE", table))
	if err != nil {
		panic(err)
	}
}

func PurgePasswords(pool *pgxpool.Pool, ctx context.Context) {
	PurgeTable(pool, ctx, "passwords")
}

func PurgeUsers(pool *pgxpool.Pool, ctx context.Context) {
	PurgeTable(pool, ctx, "users")
}

func PurgeUserViews(pool *pgxpool.Pool, ctx context.Context) {
	PurgeTable(pool, ctx, "users_view")
}

func PurgePhones(pool *pgxpool.Pool, ctx context.Context) {
	PurgeTable(pool, ctx, "phones")
}

func PurgeRoles(pool *pgxpool.Pool, ctx context.Context) {
	PurgeTable(pool, ctx, "roles")
}

func PurgeUserRoles(pool *pgxpool.Pool, ctx context.Context) {
	PurgeTable(pool, ctx, "user_roles")
}

func PurgeEmails(pool *pgxpool.Pool, ctx context.Context) {
	PurgeTable(pool, ctx, "emails")
}
