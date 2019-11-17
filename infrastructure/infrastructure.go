package infrastructure

import (
	"auth/inout"
	"auth/models"
	"auth/repositories"
	"context"
	"time"
)

type StoreInterface interface {
	// All store methods

	// Users

	CreateUser(ctx context.Context, query *inout.CreateUserRequestV1) (int, *models.User)
	GetUser(context context.Context, id int64) *models.User
	GetUsers(context context.Context, query repositories.GetUsersQuery) []*models.User

	// User Views

	GetUsersView(context context.Context, query repositories.GetUsersViewQuery) []*inout.GetUserViewResponseV1
	GetUserView(context context.Context, id int64) *inout.GetUserViewResponseV1
	CreateOrUpdateUsersView(context context.Context, query repositories.CreateOrUpdateUsersViewQuery) []*inout.GetUserViewResponseV1
	CreateOrUpdateUsersViewByUsersID(context context.Context, id []int64) []*inout.GetUserViewResponseV1
	CreateOrUpdateUsersViewByRolesID(context context.Context, id []int64) []*inout.GetUserViewResponseV1
	CreateOrUpdateUsersViewByUserID(context context.Context, id int64) []*inout.GetUserViewResponseV1
	CreateOrUpdateUsersViewByRoleID(context context.Context, id int64) []*inout.GetUserViewResponseV1
	GetUserViewFromCache(ctx context.Context, id int64) *inout.GetUserViewResponseV1
	CacheUserView(ctx context.Context, userViews []*inout.GetUserViewResponseV1)

	// Emails

	CreateEmail(ctx context.Context, userId int64, value string) (int, *models.Email)
	GetEmail(ctx context.Context, email string) (int, *models.Email)
	CreateEmailConfirmationCode(email string, code string, duration time.Duration) *models.EmailConfirmation
	GetEmailConfirmationCode(email string) string

	// Passwords

	CreatePassword(ctx context.Context, userId int64, value string) (int, *models.Password)
	GetPasswords(ctx context.Context, userId int64) []*models.Password
	GetLatestPassword(ctx context.Context, userId int64) (int, *models.Password)

	// Phones

	CreatePhone(ctx context.Context, userId int64, value string) (int, *models.Phone)
	GetPhone(ctx context.Context, phone string) (int, *models.Phone)
	CreatePhoneConfirmationCode(ctx context.Context, phone string, code string, duration time.Duration) *models.PhoneConfirmation
	GetPhoneConfirmationCode(ctx context.Context, phone string) string

	// Roles

	CreateRole(context context.Context, title string) (int, *models.Role)
	GetRole(context context.Context, id int64) (int, *models.Role)
	GetRoles(context context.Context, query repositories.GetRolesQuery) []*models.Role
}

type ESBInterface interface {
	OnUserChanged(id []int64)
	OnEmailCodeConfirmationCreated(email string, code string)
	OnPhoneCodeConfirmationCreated(phone string, code string)
	OnUsersViewChanged(usersView []*inout.GetUserViewResponseV1)
	OnPasswordChanged(userId int64)
	OnPhoneChanged(userId []int64)
	OnEmailChanged(userId []int64)
	OnRoleChanged(roleId []int64)
}

type ESBDispatcherInterface interface {
	Send(event inout.Event)
}

type PasswordProcessorInterface interface {
	Encode(context.Context, string) string
}

type AppInterface interface {
	GetStore() StoreInterface
	GetESB() ESBInterface
	GetPasswordProcessor() PasswordProcessorInterface
}
