package repositories

import (
	"auth/enums"
	"auth/functools"
	"auth/models"
	"context"
	"errors"
	"github.com/getsentry/sentry-go"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	uuid "github.com/satori/go.uuid"
	"math"
	"strings"
)

func getRolesSQL() string {
	return `
		SELECT id, created, title
		FROM roles
		WHERE (array_length($1::uuid[], 1) IS NULL OR id = ANY ($1::uuid[])) AND
		      (array_length($2::text[], 1) IS NULL OR title = ANY ($2::text[]))
		LIMIT $3;
		`
}

func createRoleSQL() string {
	return "INSERT INTO roles (title) VALUES ($1) RETURNING id, created, title;"
}

func unwrapRoleScanError(err error) int {
	var e *pgconn.PgError
	if errors.As(err, &e) && (strings.Contains(e.Message, "unique constraint \"roles_title_key\"") || strings.Contains(e.Message, "нарушает ограничение уникальности \"roles_title_key\"")) {
		return enums.RoleAlreadyExist
	} else if strings.Contains(err.Error(), "no rows") {
		return enums.RoleNotFound
	}

	sentry.CaptureException(err)
	return enums.NotOk
}

func scanRole(row pgx.Row) (int, *models.Role) {
	role := &models.Role{}

	err := row.Scan(&role.Id, &role.Created, &role.Title)
	if err != nil {
		sentry.CaptureException(err)
		return unwrapRoleScanError(err), nil
	}

	return enums.Ok, role
}

func scanRoles(rows pgx.Rows, limit int) []*models.Role {
	roles := make([]*models.Role, limit)

	var i int32

	for rows.Next() {
		_, role := scanRole(rows)
		roles[i] = role
		i++
	}

	rows.Close()

	return roles[0:i]
}

type GetRolesQuery struct {
	Pagination  *models.PaginationRequest
	Identifiers []string
	Titles      []string
}

type getRolesRawQuery struct {
	Limit       int
	Identifiers string
	Titles      string
}

func convertGetRolesQueryToRaw(query GetRolesQuery) getRolesRawQuery {

	limit := query.Pagination.Limit
	if len(query.Identifiers) > 0 {
		limit = int(math.Min(
			float64(query.Pagination.Limit),
			float64(len(query.Identifiers))))
	}

	return getRolesRawQuery{
		Limit:       limit,
		Identifiers: functools.StringsToPGArray(query.Identifiers),
		Titles:      functools.StringsToPGArray(query.Titles),
	}
}

func CreateRole(db DB, context context.Context, title string) (int, *models.Role) {
	sql := createRoleSQL()
	row := db.QueryRow(context, sql, title)
	return scanRole(row)
}

func GetRole(db DB, context context.Context, id uuid.UUID) (int, *models.Role) {
	sql := getRolesSQL()
	row := db.QueryRow(context, sql, functools.StringsToPGArray([]string{id.String()}), "{}", 1)
	return scanRole(row)
}

func GetRoles(db DB, context context.Context, query GetRolesQuery) []*models.Role {

	sql := getRolesSQL()
	rawQuery := convertGetRolesQueryToRaw(query)

	rows, err := db.Query(context, sql, rawQuery.Identifiers, rawQuery.Titles, rawQuery.Limit)
	if err != nil {
		sentry.CaptureException(err)
		return nil
	}

	return scanRoles(rows, rawQuery.Limit)
}
