package repositories

import (
	"auth/functools"
	"auth/models"
	"auth/modelsFunctools"
	"context"
	"github.com/getsentry/sentry-go"
	"github.com/jackc/pgx/v4"
	uuid "github.com/satori/go.uuid"
	"math"
)

type GetUsersViewStoreQuery struct {
	Limit      int
	Page       int
	Id         []uuid.UUID
	Roles      []uuid.UUID
	Phones     []string
	PhoneCodes []string
	Emails     []string
	EmailCodes []string
}

type CreateOrUpdateUsersViewStoreQuery struct {
	Limit int
	Id    []uuid.UUID
	Roles []uuid.UUID
}

type getUsersViewRepositoryQuery struct {
	Limit  int
	Offset int
	Id     string
	Roles  string
	Emails string
	Phones string
}

type createOrUpdateUsersViewRepositoryQuery struct {
	Limit int
	Id    string
	Roles string
}

// SQL

func getUsersViewSQL() string {
	return `
		SELECT id, created, roles, phones, emails, role_id, count(*) OVER() AS full_count
		FROM users_view u
		WHERE (array_length($1::uuid[], 1) IS NULL OR id = ANY ($1::uuid[])) AND 
		      (array_length($2::uuid[], 1) IS NULL OR ($2::uuid[]) && role_id) AND
		      (array_length($3::text[], 1) IS NULL OR ($3::text[]) && phones) AND
		      (array_length($4::text[], 1) IS NULL OR ($4::text[]) && emails)
		ORDER BY created
		LIMIT $5
		OFFSET $6;
		`
}

func updateUsersViewSQL() string {
	return `
		INSERT
		INTO users_view(id, created, roles, phones, emails, role_id)
		SELECT nuv.id, nuv.created, nuv.roles, nuv.phones, nuv.emails, nuv.role_id
		FROM users_view as cuv
				 FULL OUTER JOIN (SELECT u.id,
										 u.created,
										 array_remove(array_agg(DISTINCT r.title), NULL)::text[]          as roles,
										 array_remove(array_agg(p.value), NULL)::text[]                   as phones,
										 array_remove(array_agg(DISTINCT e.value), NULL)::text[]          as emails,
										 array_remove(array_agg(DISTINCT r.id), NULL)                     as role_id
								  FROM users u
										   LEFT JOIN emails e on u.id = e.user_id
										   LEFT JOIN phones p on u.id = p.user_id
										   LEFT JOIN user_roles ur on u.id = ur.user_id
										   LEFT JOIN roles r on ur.role_id = r.id
								  WHERE (array_length($1::uuid[], 1) IS NULL OR u.id = ANY ($1::uuid[]))
									AND (array_length($2::uuid[], 1) IS NULL OR r.id = ANY ($2::uuid[]))
								  GROUP BY u.id, ur.created, r.created, p.created, e.created) as nuv
								 ON nuv.id = cuv.id AND
									nuv.created = cuv.created AND
									nuv.phones = cuv.phones AND
									nuv.roles = cuv.roles AND
									nuv.emails = cuv.emails AND
									nuv.role_id = cuv.role_id
		WHERE cuv.id IS NULL
		ORDER BY created
		ON CONFLICT (id) DO UPDATE SET created=excluded.created,
									   phones=excluded.phones,
									   roles=excluded.roles,
									   emails=excluded.emails,
									   role_id=excluded.role_id
		RETURNING id, created, roles, phones, emails, role_id, 0;
    `
}

func scanUserView(row pgx.Row) (*models.UserView, int64) {
	userView := &models.UserView{}
	var bytes [][]byte
	var count int64

	err := row.Scan(&userView.Id, &userView.Created, &userView.Roles, &userView.Phones, &userView.Emails, &bytes, &count)
	if err != nil {
		sentry.CaptureException(err)
		return nil, count
	}

	userView.RolesID = functools.ByteArraySliceToUUIDSlice(bytes)

	return userView, count
}

func scanUsersView(rows pgx.Rows, limit int) ([]*models.UserView, int64) {

	users := make([]*models.UserView, limit)

	var i int
	var count int64

	for rows.Next() {

		user, c := scanUserView(rows)
		count = c

		if len(users) <= i {
			users = append(users, user)
		} else {
			users[i] = user
		}

		i++
	}

	rows.Close()
	return users[:i], count
}

func convertGetUsersViewQueryToRaw(query GetUsersViewStoreQuery) getUsersViewRepositoryQuery {

	maxQueryLength := functools.Max([]int{len(query.Id), len(query.Emails), len(query.Phones)})

	limit := query.Limit
	if len(query.Id) > 0 {
		limit = int(math.Min(
			float64(query.Limit),
			float64(maxQueryLength)))
	}

	limit, offset := functools.LimitPageToLimitOffset(limit, query.Page)

	return getUsersViewRepositoryQuery{
		Limit:  limit,
		Offset: offset,
		Id:     modelsFunctools.UserIDListToPGArray(query.Id),
		Roles:  modelsFunctools.RoleIDListToPGArray(query.Roles),
		Phones: functools.StringsToPGArray(query.Phones),
		Emails: functools.StringsToPGArray(query.Emails),
	}
}

func convertCreateOrUpdateUsersViewQueryToRaw(query CreateOrUpdateUsersViewStoreQuery) createOrUpdateUsersViewRepositoryQuery {

	limit := query.Limit
	if len(query.Id) > 0 {
		limit = int(math.Min(
			float64(query.Limit),
			float64(len(query.Id))))
	}

	return createOrUpdateUsersViewRepositoryQuery{
		Id:    modelsFunctools.UserIDListToPGArray(query.Id),
		Limit: limit,
		Roles: modelsFunctools.RoleIDListToPGArray(query.Roles),
	}
}

func GetUsersView(db DB, context context.Context, query GetUsersViewStoreQuery) ([]*models.UserView, *models.PaginationResponse) {
	sql := getUsersViewSQL()
	rawQuery := convertGetUsersViewQueryToRaw(query)

	rows, err := db.Query(context, sql, rawQuery.Id, rawQuery.Roles, rawQuery.Phones, rawQuery.Emails, rawQuery.Limit, rawQuery.Offset)
	if err != nil {
		sentry.CaptureException(err)
		return nil, nil
	}

	userViews, totalCount := scanUsersView(rows, rawQuery.Limit)

	return userViews, &models.PaginationResponse{
		HasNext:     functools.HasNext(totalCount, query.Limit, query.Page),
		HasPrevious: functools.HasPrevious(query.Page),
		Count:       totalCount,
	}
}

func GetUserView(db DB, context context.Context, id uuid.UUID) *models.UserView {
	sql := getUsersViewSQL()
	row := db.QueryRow(context, sql, modelsFunctools.UserIDListToPGArray([]uuid.UUID{id}), "{}", "{}", "{}", 1)
	userView, _ := scanUserView(row)
	return userView
}

func CreateOrUpdateUsersView(db DB, context context.Context, query CreateOrUpdateUsersViewStoreQuery) []*models.UserView {
	sql := updateUsersViewSQL()
	rawQuery := convertCreateOrUpdateUsersViewQueryToRaw(query)
	rows, err := db.Query(context, sql, rawQuery.Id, rawQuery.Roles)
	if err != nil {
		sentry.CaptureException(err)
		return nil
	}

	userViews, _ := scanUsersView(rows, len(query.Id))
	return userViews
}
