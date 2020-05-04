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
	"strings"
)

func createUserRoleSQL() string {
	return `
		INSERT INTO user_roles(id, user_id, role_id) 
		VALUES ($1, $2, $3) 
		RETURNING id, created, user_id, role_id, 0;
		`
}

func getUserRolesSQL() string {
	return `
		SELECT id, created, user_id, role_id, count(*) OVER() AS full_count
		FROM user_roles
		WHERE (array_length($1::uuid[], 1) IS NULL OR user_id = ANY ($1::uuid[])) AND 
		      (array_length($2::uuid[], 1) IS NULL OR role_id = ANY ($2::uuid[])) 
		LIMIT $3
		OFFSET $4;
		`
}

func deleteUserRoleSQL() string {
	return `
		DELETE FROM user_roles WHERE id = $1 RETURNING id, created, user_id, role_id, 0;
		`
}

func unwrapUserRoleScanError(err error) int {
	var e *pgconn.PgError

	if errors.As(err, &e) {
		if strings.Contains(e.Detail, "not present in table \"users\".") || strings.Contains(e.Detail, "отсутствует в таблице \"users\"") {
			return enums.UserNotFound
		} else if strings.Contains(e.Detail, "not present in table \"roles\".") || strings.Contains(e.Detail, "отсутствует в таблице \"roles\"") {
			return enums.RoleNotFound
		} else if strings.Contains(e.Message, "violates unique constraint \"user_roles_pkey\"") {
			return enums.UserRoleAlreadyExist
		} else if strings.Contains(e.Message, "duplicate key value violates unique constraint \"user_roles_idx\"") || strings.Contains(e.Message, "повторяющееся значение ключа нарушает ограничение уникальности \"user_roles_idx\"") {
			return enums.UserRoleAlreadyExist
		}
	} else if strings.Contains(err.Error(), "no rows in result") {
		return enums.UserRoleNotFound
	}

	sentry.CaptureException(err)
	return enums.NotOk
}

func scanUserRole(row pgx.Row) (int, *models.UserRole, int64) {
	ur := &models.UserRole{}
	var count int64

	err := row.Scan(&ur.Id, &ur.Created, &ur.UserId, &ur.RoleId, &count)
	if err != nil {
		sentry.CaptureException(err)
		return unwrapUserRoleScanError(err), nil, 0
	}

	return enums.Ok, ur, count
}

func scanUserRoles(rows pgx.Rows, limit int) ([]*models.UserRole, int64) {
	userRoles := make([]*models.UserRole, limit)

	var i int32
	var count int64

	for rows.Next() {
		_, ur, c := scanUserRole(rows)
		count = c
		userRoles[i] = ur
		i++
	}

	rows.Close()

	return userRoles[0:i], count
}

type GetUserRoleQuery struct {
	UserId     []uuid.UUID
	RoleId     []uuid.UUID
	Pagination *models.PaginationRequest
}

type getUserRoleRawQuery struct {
	Offset int
	Limit  int
	UserId string
	RoleId string
}

func convertGetUserRoleQueryToRaw(query GetUserRoleQuery) getUserRoleRawQuery {
	limit, offset := functools.LimitPageToLimitOffset(query.Pagination.Limit, query.Pagination.Page)
	return getUserRoleRawQuery{
		Limit:  limit,
		Offset: offset,
		UserId: functools.UUIDListToPGArray(query.UserId),
		RoleId: functools.UUIDListToPGArray(query.RoleId),
	}
}

func CreateUserRole(db DB, ctx context.Context, userID uuid.UUID, roleID uuid.UUID) (int, *models.UserRole) {
	sql := createUserRoleSQL()
	row := db.QueryRow(ctx, sql, uuid.NewV4(), userID, roleID)
	status, userRole, _ := scanUserRole(row)
	return status, userRole
}

func GetUserRoles(db DB, ctx context.Context, query GetUserRoleQuery) ([]*models.UserRole, *models.PaginationResponse) {
	sql := getUserRolesSQL()
	rawQuery := convertGetUserRoleQueryToRaw(query)
	rows, err := db.Query(ctx, sql, rawQuery.UserId, rawQuery.RoleId, rawQuery.Limit, rawQuery.Offset)
	if err != nil {
		sentry.CaptureException(err)
		return nil, nil
	}

	userRoles, totalCount := scanUserRoles(rows, rawQuery.Limit)

	return userRoles, &models.PaginationResponse{
		HasNext:     functools.HasNext(totalCount, query.Pagination.Limit, query.Pagination.Page),
		HasPrevious: functools.HasPrevious(query.Pagination.Page),
		Count:       totalCount,
	}
}

func DeleteUserRole(db DB, ctx context.Context, id uuid.UUID) (int, *models.UserRole) {
	sql := deleteUserRoleSQL()
	row := db.QueryRow(ctx, sql, id)
	status, userRole, _ := scanUserRole(row)
	return status, userRole
}
