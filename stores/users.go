package stores

import (
	"auth/enums"
	"auth/functools"
	"auth/inout"
	"auth/models"
	"auth/repositories"
	"context"
	"github.com/getsentry/sentry-go"
	uuid "github.com/satori/go.uuid"
)

func createPhoneForUser(tx repositories.DB, ctx context.Context, phone string, userId uuid.UUID) int {
	status, _ := repositories.CreatePhone(tx, ctx, userId, phone)
	return status
}

func createEmailForUser(tx repositories.DB, ctx context.Context, email string, userId uuid.UUID) int {
	status, _ := repositories.CreateEmail(tx, ctx, userId, email)
	return status
}

func createPasswordForUser(tx repositories.DB, ctx context.Context, password string, userId uuid.UUID) int {
	status, _ := repositories.CreatePassword(tx, ctx, userId, password)
	return status
}

func (store *DatabaseStore) CreateUser(ctx context.Context, query *inout.CreateUserResponseV1_Request) (int, *models.User) {
	tx, err := store.db.Begin(ctx)

	if tx == nil {
		return enums.NotOk, nil
	}

	if repositories.Rollback(tx, ctx, err != nil) {
		sentry.CaptureException(err)
		return enums.NotOk, nil
	}

	user := repositories.CreateUser(tx, ctx)

	if repositories.Rollback(tx, ctx, user == nil) {
		return enums.NotOk, nil
	}

	var statuses []int

	if query.Phone != "" {
		statuses = append(statuses, createPhoneForUser(tx, ctx, query.Phone, user.Id))
	}

	if query.Email != "" {
		statuses = append(statuses, createEmailForUser(tx, ctx, query.Email, user.Id))
	}

	if query.Password != "" {
		statuses = append(statuses, createPasswordForUser(tx, ctx, query.Password, user.Id))
	}

	if repositories.Rollback(tx, ctx, !functools.All(enums.Ok, statuses)) {
		return functools.Max(statuses), nil
	}

	err = tx.Commit(ctx)

	if err != nil {
		sentry.CaptureException(err)
		return enums.NotOk, nil
	}

	return enums.Ok, user
}

func (store *DatabaseStore) GetUser(context context.Context, id uuid.UUID) *models.User {
	return repositories.GetUser(store.db, context, id)
}

func (store *DatabaseStore) GetUsers(context context.Context, query repositories.GetUsersQuery) []*models.User {
	return repositories.GetUsers(store.db, context, query)
}

func (store *DatabaseStore) DeleteUser(ctx context.Context, id uuid.UUID) (int, *models.User) {
	return repositories.DeleteUser(store.db, ctx, id)
}
