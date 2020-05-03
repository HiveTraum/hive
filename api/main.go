package api

import (
	"auth/auth"
	"auth/config"
	"auth/controllers"
	"auth/presenters"
	"github.com/getsentry/sentry-go"
	"github.com/golang/mock/gomock"
	"net/http"
	"strconv"
)

type IAPI interface {
	GetAuthenticationController() auth.IAuthenticationController
	GetController() controllers.IController

	// Emails

	CreateEmailV1(w http.ResponseWriter, r *http.Request)
	CreateEmailConfirmationV1(w http.ResponseWriter, r *http.Request)

	// Passwords

	CreatePasswordV1(w http.ResponseWriter, r *http.Request)

	// Phones

	CreatePhoneV1(w http.ResponseWriter, r *http.Request)
	CreatePhoneConfirmationV1(w http.ResponseWriter, r *http.Request)

	// Roles

	CreateRoleV1(w http.ResponseWriter, r *http.Request)
	GetRolesV1(w http.ResponseWriter, r *http.Request)
	GetRoleV1(w http.ResponseWriter, r *http.Request)

	// Secrets

	GetSecretV1(w http.ResponseWriter, r *http.Request)

	// Sessions

	CreateSessionV1(w http.ResponseWriter, r *http.Request)

	// Users

	GetUserV1(w http.ResponseWriter, r *http.Request)
	DeleteUserV1(w http.ResponseWriter, r *http.Request)
	GetUsersV1(w http.ResponseWriter, r *http.Request)
	CreateUserV1(w http.ResponseWriter, r *http.Request)

	// User Views

	GetUserViewV1(w http.ResponseWriter, r *http.Request)
	GetUsersViewV1(w http.ResponseWriter, r *http.Request)

	// User Roles

	DeleteUserRoleV1(w http.ResponseWriter, r *http.Request)
	CreateUserRoleV1(w http.ResponseWriter, r *http.Request)
	GetUserRolesV1(w http.ResponseWriter, r *http.Request)
}

type API struct {
	Controller               controllers.IController
	authenticationController auth.IAuthenticationController
	Parser                   *presenters.Parser
	Renderer                 *presenters.Renderer
	environment              *config.Environment
}

func (api *API) GetLoginController() auth.IAuthenticationController {
	return api.authenticationController
}

func InitAPI(controller controllers.IController, authenticationController auth.IAuthenticationController, environment *config.Environment) *API {
	return &API{
		Controller:               controller,
		authenticationController: authenticationController,
		Parser:                   presenters.InitParser(),
		Renderer:                 presenters.InitRenderer(),
		environment:              environment,
	}
}

type APIWithMockedInternals struct {
	API                      *API
	Controller               *controllers.MockIController
	AuthenticationController *auth.MockIAuthenticationController
}

func InitAPIWithMockedInternals(ctrl *gomock.Controller) *APIWithMockedInternals {
	controller := controllers.NewMockIController(ctrl)
	authenticationController := auth.NewMockIAuthenticationController(ctrl)
	return &APIWithMockedInternals{
		API:                      InitAPI(controller, authenticationController, config.InitEnvironment()),
		Controller:               controller,
		AuthenticationController: authenticationController,
	}
}

func unhandledStatus(r *http.Request, status int) int {

	request := sentry.Request{}
	request.FromHTTPRequest(r)

	sentry.CaptureEvent(&sentry.Event{
		Level:   sentry.LevelError,
		Message: "Unhandled controller status",
		Tags:    map[string]string{"controller status": strconv.Itoa(status)},
		Request: request,
	})

	return http.StatusInternalServerError
}
