package repositories

import (
	"context"
	"github.com/getsentry/sentry-go"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type DB interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)
}

func Rollback(tx pgx.Tx, ctx context.Context, condition bool) bool {

	if condition == false {
		return false
	}

	err := tx.Rollback(ctx)
	if err != nil {
		sentry.CaptureException(err)
	}

	return condition
}
