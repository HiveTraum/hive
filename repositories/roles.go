package repositories

import (
	"auth/functools"
	"auth/models"
	"context"
	"github.com/getsentry/sentry-go"
	"github.com/jackc/pgx/v4"
	"math"
)

func getRolesSQL() string {
	return `
		SELECT id, created, title
		FROM roles
		WHERE (array_length($1::integer[], 1) IS NULL OR id = ANY ($1::bigint[]))
		LIMIT $2;
		`
}

func createRoleSQL() string {
	return "INSERT INTO roles (title) VALUES ($1) RETURNING id, created, title;"
}

func scanRole(row pgx.Row) *models.Role {
	role := &models.Role{}

	err := row.Scan(&role.Id, &role.Created, &role.Title)
	if err != nil {
		sentry.CaptureException(err)
		return nil
	}

	return role
}

func scanRoles(rows pgx.Rows, limit int) []*models.Role {
	roles := make([]*models.Role, limit)

	var i int32

	for rows.Next() {
		roles[i] = scanRole(rows)
		i++
	}

	rows.Close()

	return roles[0:i]
}

type GetRolesQuery struct {
	Limit       int
	Identifiers []int64
}

type getRolesRawQuery struct {
	Limit       int
	Identifiers string
}

func convertGetRolesQueryToRaw(query GetRolesQuery) getRolesRawQuery {

	limit := query.Limit
	if len(query.Identifiers) > 0 {
		limit = int(math.Min(
			float64(query.Limit),
			float64(len(query.Identifiers))))
	}

	return getRolesRawQuery{
		Limit:       limit,
		Identifiers: functools.Int64ListToPGArray(query.Identifiers),
	}
}

func CreateRole(db DB, context context.Context, title string) *models.Role {
	sql := createRoleSQL()
	row := db.QueryRow(context, sql, title)
	return scanRole(row)
}

func GetRole(db DB, context context.Context, id int64) *models.Role {
	sql := getRolesSQL()
	row := db.QueryRow(context, sql, functools.Int64ListToPGArray([]int64{id}), 1)
	return scanRole(row)
}

func GetRoles(db DB, context context.Context, query GetRolesQuery) []*models.Role {

	sql := getRolesSQL()
	rawQuery := convertGetRolesQueryToRaw(query)

	rows, err := db.Query(context, sql, rawQuery.Identifiers, rawQuery.Limit)
	if err != nil {
		sentry.CaptureException(err)
		return nil
	}

	return scanRoles(rows, rawQuery.Limit)
}
