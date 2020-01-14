package infrastructure

import (
	"auth/inout"
	"auth/models"
	"auth/repositories"
	"context"
	"time"
)

type StoreInterface interface {
	// Store can be used to combine multiple physical storage elements, like postgres, redis, elasticSearch and etc...

	// All store methods

	// Users

	CreateUser(ctx context.Context, query *inout.CreateUserRequestV1) (int, *models.User)
	GetUser(context context.Context, id models.UserID) *models.User
	GetUsers(context context.Context, query repositories.GetUsersQuery) []*models.User

	// User Views

	GetUsersView(context context.Context, query repositories.GetUsersViewStoreQuery) []*inout.GetUserViewResponseV1
	GetUserView(context context.Context, id models.UserID) *inout.GetUserViewResponseV1
	CreateOrUpdateUsersView(context context.Context, query repositories.CreateOrUpdateUsersViewStoreQuery) []*inout.GetUserViewResponseV1
	CreateOrUpdateUsersViewByUsersID(context context.Context, id []models.UserID) []*inout.GetUserViewResponseV1
	CreateOrUpdateUsersViewByRolesID(context context.Context, id []models.RoleID) []*inout.GetUserViewResponseV1
	CreateOrUpdateUsersViewByUserID(context context.Context, id models.UserID) []*inout.GetUserViewResponseV1
	CreateOrUpdateUsersViewByRoleID(context context.Context, id models.RoleID) []*inout.GetUserViewResponseV1
	CacheUserView(ctx context.Context, userViews []*inout.GetUserViewResponseV1)

	// Emails

	CreateEmail(ctx context.Context, userId models.UserID, value string) (int, *models.Email)
	GetEmail(ctx context.Context, email string) (int, *models.Email)
	CreateEmailConfirmationCode(ctx context.Context, email string, code string, duration time.Duration) *models.EmailConfirmation
	GetEmailConfirmationCode(ctx context.Context, email string) string

	// Passwords

	CreatePassword(ctx context.Context, userId models.UserID, value string) (int, *models.Password)
	GetPasswords(ctx context.Context, userId models.UserID) []*models.Password
	GetLatestPassword(ctx context.Context, userId models.UserID) (int, *models.Password)

	// Phones

	CreatePhone(ctx context.Context, userId models.UserID, value string) (int, *models.Phone)
	GetPhone(ctx context.Context, phone string) (int, *models.Phone)
	CreatePhoneConfirmationCode(ctx context.Context, phone string, code string, duration time.Duration) *models.PhoneConfirmation
	GetPhoneConfirmationCode(ctx context.Context, phone string) string

	// Roles

	CreateRole(context context.Context, title string) (int, *models.Role)
	GetRole(context context.Context, id models.RoleID) (int, *models.Role)
	GetRoles(context context.Context, query repositories.GetRolesQuery) []*models.Role

	// User Roles

	CreateUserRole(ctx context.Context, userId models.UserID, roleId models.RoleID) (int, *models.UserRole)
	GetUserRoles(ctx context.Context, query repositories.GetUserRoleQuery) []*models.UserRole
	DeleteUserRole(ctx context.Context, id models.UserRoleID) (int, *models.UserRole)

	// Secrets

	GetSecret(ctx context.Context, id models.SecretID) *models.Secret
	GetActualSecret(ctx context.Context) *models.Secret

	// Sessions

	CreateSession(ctx context.Context, fingerprint string, userID models.UserID, secretID models.SecretID, userAgent string) (int, *models.Session)
	GetSession(ctx context.Context, fingerprint string, refreshToken string, userID models.UserID) *models.Session
}

type ESBInterface interface {
	OnUserChanged(id []models.UserID)
	OnEmailCodeConfirmationCreated(email string, code string)
	OnPhoneCodeConfirmationCreated(phone string, code string)
	OnUsersViewChanged(usersView []*inout.GetUserViewResponseV1)
	OnPasswordChanged(userId models.UserID)
	OnPhoneChanged(userId []models.UserID)
	OnEmailChanged(userId []models.UserID)
	OnRoleChanged(roleId []models.RoleID)
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

	DecodeAccessToken(ctx context.Context, token string, secret string) (int, *models.AccessTokenPayload)
	DecodeAccessTokenWithoutValidation(ctx context.Context, tokenValue string) (int, *models.AccessTokenPayload)
	EncodeAccessToken(ctx context.Context, userID models.UserID, roles []string, secret string, expire time.Time) string

	EncodePassword(context.Context, string) string
	VerifyPassword(ctx context.Context, password string, encodedPassword string) bool
}

type AppInterface interface {
	GetStore() StoreInterface
	GetESB() ESBInterface
	GetLoginController() LoginControllerInterface
}
