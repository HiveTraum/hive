package repositories

import (
	"hive/enums"
	"hive/functools"
	"hive/models"
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
		SELECT id, created, title, count(*) OVER() AS full_count
		FROM roles
		WHERE (array_length($1::uuid[], 1) IS NULL OR id = ANY ($1::uuid[])) AND
		      (array_length($2::text[], 1) IS NULL OR title = ANY ($2::text[]))
		LIMIT $3 
		OFFSET $4;
		`
}

func createRoleSQL() string {
	return "INSERT INTO roles (id, title) VALUES ($1, $2) RETURNING id, created, title, 0;"
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

func scanRole(row pgx.Row) (int, *models.Role, int64) {
	role := &models.Role{}
	var count int64

	err := row.Scan(&role.Id, &role.Created, &role.Title, &count)
	if err != nil {
		sentry.CaptureException(err)
		return unwrapRoleScanError(err), nil, 0
	}

	return enums.Ok, role, count
}

func scanRoles(rows pgx.Rows, limit int) ([]*models.Role, int64) {
	roles := make([]*models.Role, limit)

	var i int32
	var count int64

	for rows.Next() {
		_, role, c := scanRole(rows)
		roles[i] = role
		count = c
		i++
	}

	rows.Close()

	return roles[0:i], count
}

type GetRolesQuery struct {
	Pagination  *models.PaginationRequest
	Identifiers []string
	Titles      []string
}

type getRolesRawQuery struct {
	Limit       int
	Offset      int
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

	limit, offset := functools.LimitPageToLimitOffset(limit, query.Pagination.Page)

	return getRolesRawQuery{
		Limit:       limit,
		Offset:      offset,
		Identifiers: functools.StringsToPGArray(query.Identifiers),
		Titles:      functools.StringsToPGArray(query.Titles),
	}
}

func CreateRole(db DB, context context.Context, title string) (int, *models.Role) {
	sql := createRoleSQL()
	row := db.QueryRow(context, sql, uuid.NewV4(), title)
	status, role, _ := scanRole(row)
	return status, role
}

func GetRole(db DB, context context.Context, id uuid.UUID) (int, *models.Role) {
	sql := getRolesSQL()
	row := db.QueryRow(context, sql, functools.StringsToPGArray([]string{id.String()}), "{}", 1, 0)
	status, role, _ := scanRole(row)
	return status, role
}

func GetRoles(db DB, context context.Context, query GetRolesQuery) ([]*models.Role, *models.PaginationResponse) {

	sql := getRolesSQL()
	rawQuery := convertGetRolesQueryToRaw(query)

	rows, err := db.Query(context, sql, rawQuery.Identifiers, rawQuery.Titles, rawQuery.Limit, rawQuery.Offset)
	if err != nil {
		sentry.CaptureException(err)
		return nil, nil
	}

	roles, count := scanRoles(rows, rawQuery.Limit)

	return roles, &models.PaginationResponse{
		HasNext:     functools.HasNext(count, query.Pagination.Limit, query.Pagination.Page),
		HasPrevious: functools.HasPrevious(query.Pagination.Page),
		Count:       count,
	}
}
