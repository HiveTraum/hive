package repositories

import (
	"auth/enums"
	"auth/models"
	"auth/modelsFunctools"
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
		INSERT INTO user_roles(user_id, role_id) 
		VALUES ($1, $2) 
		RETURNING id, created, user_id, role_id;
		`
}

func getUserRolesSQL() string {
	return `
		SELECT id, created, user_id, role_id 
		FROM user_roles
		WHERE (array_length($1::uuid[], 1) IS NULL OR user_id = ANY ($1::uuid[])) AND 
		      (array_length($2::uuid[], 1) IS NULL OR role_id = ANY ($2::uuid[])) 
		LIMIT $3;
		`
}

func deleteUserRoleSQL() string {
	return `
		DELETE FROM user_roles WHERE id = $1 RETURNING id, created, user_id, role_id;
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

func scanUserRole(row pgx.Row) (int, *models.UserRole) {
	ur := &models.UserRole{}

	err := row.Scan(&ur.Id, &ur.Created, &ur.UserId, &ur.RoleId)
	if err != nil {
		sentry.CaptureException(err)
		return unwrapUserRoleScanError(err), nil
	}

	return enums.Ok, ur
}

func scanUserRoles(rows pgx.Rows, limit int) []*models.UserRole {
	userRoles := make([]*models.UserRole, limit)

	var i int32

	for rows.Next() {
		_, ur := scanUserRole(rows)
		userRoles[i] = ur
		i++
	}

	rows.Close()

	return userRoles[0:i]
}

type GetUserRoleQuery struct {
	UserId     []uuid.UUID
	RoleId     []uuid.UUID
	Pagination *models.PaginationRequest
}

type getUserRoleRawQuery struct {
	Limit  int
	UserId string
	RoleId string
}

func convertGetUserRoleQueryToRaw(query GetUserRoleQuery) getUserRoleRawQuery {
	return getUserRoleRawQuery{
		Limit:  query.Pagination.Limit,
		UserId: modelsFunctools.UserIDListToPGArray(query.UserId),
		RoleId: modelsFunctools.RoleIDListToPGArray(query.RoleId),
	}
}

func CreateUserRole(db DB, ctx context.Context, userID uuid.UUID, roleID uuid.UUID) (int, *models.UserRole) {
	sql := createUserRoleSQL()
	row := db.QueryRow(ctx, sql, userID, roleID)
	return scanUserRole(row)
}

func GetUserRoles(db DB, ctx context.Context, query GetUserRoleQuery) []*models.UserRole {
	sql := getUserRolesSQL()
	rawQuery := convertGetUserRoleQueryToRaw(query)
	rows, err := db.Query(ctx, sql, rawQuery.UserId, rawQuery.RoleId, rawQuery.Limit)
	if err != nil {
		sentry.CaptureException(err)
		return nil
	}

	return scanUserRoles(rows, rawQuery.Limit)
}

func DeleteUserRole(db DB, ctx context.Context, id uuid.UUID) (int, *models.UserRole) {
	sql := deleteUserRoleSQL()
	row := db.QueryRow(ctx, sql, id)
	return scanUserRole(row)
}
