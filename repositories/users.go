package repositories

import (
	"auth/enums"
	"auth/models"
	"auth/modelsFunctools"
	"context"
	"github.com/getsentry/sentry-go"
	"github.com/jackc/pgx/v4"
	uuid "github.com/satori/go.uuid"
	"math"
)

func deleteUserSQL() string {
	return "DELETE FROM users WHERE id = $1 RETURNING id, created"
}

func createUserSQL() string {
	return "INSERT INTO users (id, created) VALUES ($1, default) RETURNING id, created;"
}

func getUsersSQL() string {
	return `
		SELECT id, created
		FROM users
		WHERE (array_length($1::uuid[], 1) IS NULL OR id = ANY ($1::uuid[]))
		LIMIT $2;
		`
}

func scanUser(row pgx.Row) *models.User {
	user := &models.User{}

	err := row.Scan(&user.Id, &user.Created)
	if err != nil {
		sentry.CaptureException(err)
		return nil
	}

	return user
}

func scanUsers(rows pgx.Rows, limit int) []*models.User {
	users := make([]*models.User, limit)

	var i int32

	for rows.Next() {
		users[i] = scanUser(rows)
		i++
	}

	rows.Close()

	return users[0:i]
}

type GetUsersQuery struct {
	Limit int
	Page  int
	Id    []uuid.UUID
}

type getUsersRawQuery struct {
	Limit int
	Id    string
}

func convertGetUsersQueryToRaw(query GetUsersQuery) getUsersRawQuery {

	limit := query.Limit
	if len(query.Id) > 0 {
		limit = int(math.Min(
			float64(query.Limit),
			float64(len(query.Id))))
	}

	return getUsersRawQuery{
		Limit: limit,
		Id:    modelsFunctools.UserIDListToPGArray(query.Id),
	}
}

func CreateUser(db DB, ctx context.Context) *models.User {
	sql := createUserSQL()
	row := db.QueryRow(ctx, sql, uuid.NewV4())
	return scanUser(row)
}

func GetUser(db DB, context context.Context, id uuid.UUID) *models.User {
	sql := getUsersSQL()
	row := db.QueryRow(context, sql, modelsFunctools.UserIDListToPGArray([]uuid.UUID{id}), 1)
	return scanUser(row)
}

func GetUsers(db DB, context context.Context, query GetUsersQuery) []*models.User {

	sql := getUsersSQL()
	rawQuery := convertGetUsersQueryToRaw(query)
	rows, err := db.Query(context, sql, rawQuery.Id, rawQuery.Limit)
	if err != nil {
		sentry.CaptureException(err)
		return nil
	}

	return scanUsers(rows, rawQuery.Limit)
}

func DeleteUser(db DB, ctx context.Context, id uuid.UUID) (int, *models.User) {
	sql := deleteUserSQL()
	row := db.QueryRow(ctx, sql, id)
	deletedUser := scanUser(row)
	if deletedUser == nil {
		return enums.UserNotFound, nil
	} else {
		return enums.Ok, deletedUser
	}
}
