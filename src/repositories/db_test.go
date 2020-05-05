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

func PurgeUserRoles(pool *pgxpool.Pool, ctx context.Context) {
	PurgeTable(pool, ctx, "user_roles")
}

func PurgePasswords(pool *pgxpool.Pool, ctx context.Context) {
	PurgeTable(pool, ctx, "passwords")
}

func SetSeqScan(pool *pgxpool.Pool, ctx context.Context, onOff bool) {
	sql := fmt.Sprintf("SET enable_seqscan = %t", onOff)
	_, err := pool.Exec(ctx, sql)
	if err != nil {
		panic(err)
	}
}

func Explain(pool *pgxpool.Pool, ctx context.Context, sql string, args ...interface{}) []string {

	sqlWithExplain := fmt.Sprintf("EXPLAIN ANALYZE %s", sql)

	rows, err := pool.Query(ctx, sqlWithExplain, args...)
	if err != nil {
		panic(err)
	}

	var strings []string

	for rows.Next() {
		var s string
		err := rows.Scan(&s)

		if err != nil {
			panic(err)
		}

		strings = append(strings, s)
	}

	return strings
}
