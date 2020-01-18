package infrastructure

import (
	"auth/inout"
	"auth/models"
	"auth/repositories"
	"context"
	uuid "github.com/satori/go.uuid"
	"time"
)

type StoreInterface interface {
	// Store can be used to combine multiple physical storage elements, like postgres, redis, elasticSearch and etc...

	// All store methods

	// Users

	CreateUser(ctx context.Context, query *inout.CreateUserRequestV1) (int, *models.User)
	GetUser(context context.Context, id uuid.UUID) *models.User
	GetUsers(context context.Context, query repositories.GetUsersQuery) []*models.User

	// User Views

	GetUsersView(context context.Context, query repositories.GetUsersViewStoreQuery) []*models.UserView
	GetUserView(context context.Context, id uuid.UUID) *models.UserView
	CreateOrUpdateUsersView(context context.Context, query repositories.CreateOrUpdateUsersViewStoreQuery) []*models.UserView
	CreateOrUpdateUsersViewByUsersID(context context.Context, id []uuid.UUID) []*models.UserView
	CreateOrUpdateUsersViewByRolesID(context context.Context, id []uuid.UUID) []*models.UserView
	CreateOrUpdateUsersViewByUserID(context context.Context, id uuid.UUID) []*models.UserView
	CreateOrUpdateUsersViewByRoleID(context context.Context, id uuid.UUID) []*models.UserView
	CacheUserView(ctx context.Context, userViews []*models.UserView)

	// Emails

	CreateEmail(ctx context.Context, userId uuid.UUID, value string) (int, *models.Email)
	GetEmail(ctx context.Context, email string) (int, *models.Email)
	CreateEmailConfirmationCode(ctx context.Context, email string, code string, duration time.Duration) *models.EmailConfirmation
	GetEmailConfirmationCode(ctx context.Context, email string) string

	// Passwords

	CreatePassword(ctx context.Context, userId uuid.UUID, value string) (int, *models.Password)
	GetPasswords(ctx context.Context, userId uuid.UUID) []*models.Password
	GetLatestPassword(ctx context.Context, userId uuid.UUID) (int, *models.Password)

	// Phones

	CreatePhone(ctx context.Context, userId uuid.UUID, value string) (int, *models.Phone)
	GetPhone(ctx context.Context, phone string) (int, *models.Phone)
	CreatePhoneConfirmationCode(ctx context.Context, phone string, code string, duration time.Duration) *models.PhoneConfirmation
	GetPhoneConfirmationCode(ctx context.Context, phone string) string

	// Roles

	CreateRole(context context.Context, title string) (int, *models.Role)
	GetRole(context context.Context, id uuid.UUID) (int, *models.Role)
	GetRoles(context context.Context, query repositories.GetRolesQuery) []*models.Role

	// User Roles

	CreateUserRole(ctx context.Context, userId uuid.UUID, roleId uuid.UUID) (int, *models.UserRole)
	GetUserRoles(ctx context.Context, query repositories.GetUserRoleQuery) []*models.UserRole
	DeleteUserRole(ctx context.Context, id uuid.UUID) (int, *models.UserRole)

	// Secrets

	GetSecret(ctx context.Context, id uuid.UUID) *models.Secret
	GetActualSecret(ctx context.Context) *models.Secret

	// Sessions

	CreateSession(ctx context.Context, fingerprint string, userID uuid.UUID, secretID uuid.UUID, userAgent string) (int, *models.Session)
	GetSession(ctx context.Context, fingerprint string, refreshToken string, userID uuid.UUID) *models.Session
}

type ESBInterface interface {
	OnUserChanged(id []uuid.UUID)
	OnEmailCodeConfirmationCreated(email string, code string)
	OnPhoneCodeConfirmationCreated(phone string, code string)
	OnUsersViewChanged(usersView []*models.UserView)
	OnPasswordChanged(userId uuid.UUID)
	OnPhoneChanged(userId []uuid.UUID)
	OnEmailChanged(userId []uuid.UUID)
	OnRoleChanged(roleId []uuid.UUID)
}

type ESBDispatcherInterface interface {
	Send(event inout.Event)
}

type LoginControllerInterface interface {
	Login(ctx context.Context, credentials inout.CreateSessionRequestV1) (int, *models.User)

	LoginByTokens(ctx context.Context, refreshToken string, accessToken string, fingerprint string) (int, *models.User)
	LoginByEmailAndCode(ctx context.Context, email string, emailCode string) (int, *models.User)
	LoginByEmailAndPassword(ctx context.Context, email string, password string) (int, *models.User)
	LoginByPhoneAndCode(ctx context.Context, phone string, phoneCode string) (int, *models.User)
	LoginByPhoneAndPassword(ctx context.Context, phone string, password string) (int, *models.User)

	NormalizePhone(ctx context.Context, phone string) string
	NormalizeEmail(ctx context.Context, email string) string

	DecodeAccessToken(ctx context.Context, token string, secret uuid.UUID) (int, *models.AccessTokenPayload)
	DecodeAccessTokenWithoutValidation(ctx context.Context, tokenValue string) (int, *models.AccessTokenPayload)
	EncodeAccessToken(ctx context.Context, userID uuid.UUID, roles []string, secret uuid.UUID, expire time.Time) string

	EncodePassword(context.Context, string) string
	VerifyPassword(ctx context.Context, password string, encodedPassword string) bool
}

type AppInterface interface {
	GetStore() StoreInterface
	GetESB() ESBInterface
	GetLoginController() LoginControllerInterface
}
