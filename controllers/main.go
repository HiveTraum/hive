package controllers

import (
	"auth/config"
	"auth/eventDispatchers"
	"auth/models"
	"auth/passwordProcessors"
	"auth/repositories"
	"auth/stores"
	"context"
	"github.com/golang/mock/gomock"
	uuid "github.com/satori/go.uuid"
)

type IController interface {

	// Secrets

	GetActualSecret(ctx context.Context) *models.Secret
	GetSecret(ctx context.Context, id uuid.UUID) *models.Secret

	// Sessions

	CreateSession(ctx context.Context, userID uuid.UUID, userAgent string, fingerprint string) (int, *models.Session)

	// Passwords

	CreatePassword(ctx context.Context, userId uuid.UUID, value string) (int, *models.Password)

	// Emails

	CreateEmailConfirmation(ctx context.Context, email string) (int, *models.EmailConfirmation)
	CreateEmail(ctx context.Context, email string, code string, userId uuid.UUID) (int, *models.Email)

	// Users

	CreateUser(ctx context.Context, password, email, emailCode, phone, phoneCode string) (int, *models.User)
	GetUsers(ctx context.Context, query repositories.GetUsersQuery) []*models.User
	DeleteUser(ctx context.Context, id uuid.UUID) (int, *models.User)
	GetUser(ctx context.Context, id uuid.UUID) *models.User

	// Phone

	CreatePhoneConfirmation(ctx context.Context, phone string) (int, *models.PhoneConfirmation)
	CreatePhone(ctx context.Context, phone string, code string, userId uuid.UUID) (int, *models.Phone)

	// Roles

	GetRole(ctx context.Context, id uuid.UUID) (int, *models.Role)
	CreateRole(ctx context.Context, title string) (int, *models.Role)
	GetRoles(ctx context.Context, query repositories.GetRolesQuery) ([]*models.Role, *models.PaginationResponse)

	// User Roles

	GetUserRoles(ctx context.Context, query repositories.GetUserRoleQuery) ([]*models.UserRole, *models.PaginationResponse)
	CreateUserRole(ctx context.Context, userId uuid.UUID, roleID uuid.UUID) (int, *models.UserRole)
	DeleteUserRole(ctx context.Context, id uuid.UUID) (int, *models.UserRole)

	// User Views

	CreateOrUpdateUsersView(ctx context.Context, id []uuid.UUID) []*models.UserView
	CreateOrUpdateUsersViewByRoles(ctx context.Context, rolesIds []uuid.UUID) []*models.UserView
	GetUserView(ctx context.Context, id uuid.UUID) *models.UserView
	GetUserViews(ctx context.Context, query repositories.GetUsersViewStoreQuery) ([]*models.UserView, *models.PaginationResponse)

	// Events

	OnUserChanged(id []uuid.UUID)
	OnEmailCodeConfirmationCreated(email string, code string)
	OnPhoneCodeConfirmationCreated(phone string, code string)
	OnUsersViewChanged(usersView []*models.UserView)
	OnPasswordChanged(userId uuid.UUID)
	OnPhoneChanged(userId []uuid.UUID)
	OnEmailChanged(userId []uuid.UUID)
	OnRoleChanged(roleId []uuid.UUID)
	OnSecretCreatedV1(secret *models.Secret)
}

type Controller struct {
	store             stores.IStore
	passwordProcessor passwordProcessors.IPasswordProcessor
	dispatcher        eventDispatchers.IEventDispatcher
	environment       *config.Environment
}

func InitController(store stores.IStore, passwordProcessor passwordProcessors.IPasswordProcessor, dispatcher eventDispatchers.IEventDispatcher, environment *config.Environment) *Controller {
	return &Controller{
		store:             store,
		passwordProcessor: passwordProcessor,
		dispatcher:        dispatcher,
		environment:       environment,
	}
}

type ControllerWithMockedInternals struct {
	Controller        *Controller
	Dispatcher        *eventDispatchers.MockIEventDispatcher
	Store             *stores.MockIStore
	PasswordProcessor *passwordProcessors.MockIPasswordProcessor
}

func InitControllerWithMockedInternals(ctrl *gomock.Controller) *ControllerWithMockedInternals {
	dispatcher := eventDispatchers.NewMockIEventDispatcher(ctrl)
	store := stores.NewMockIStore(ctrl)
	passwordProcessor := passwordProcessors.NewMockIPasswordProcessor(ctrl)
	environment := config.InitEnvironment()
	return &ControllerWithMockedInternals{
		Controller:        InitController(store, passwordProcessor, dispatcher, environment),
		Dispatcher:        dispatcher,
		Store:             store,
		PasswordProcessor: passwordProcessor,
	}
}
