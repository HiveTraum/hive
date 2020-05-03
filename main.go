package main

import (
	api2 "auth/api"
	"auth/auth"
	"auth/auth/backends"
	"auth/config"
	"auth/controllers"
	"auth/enums"
	"auth/eventDispatchers"
	"auth/middlewares"
	"auth/passwordProcessors"
	"auth/repositories/inMemoryRepository"
	"auth/repositories/postgresRepository"
	"auth/repositories/redisRepository"
	"auth/stores"
	"context"
	"fmt"
	sentryHttp "github.com/getsentry/sentry-go/http"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
	"time"
)

func InitialAdmin(environment *config.Environment, store stores.IStore) {
	if environment.InitialAdmin == "" {
		return
	}

	ctx := context.Background()

	emailAndPassword := strings.Split(environment.InitialAdmin, ":")
	emailValue := emailAndPassword[0]
	passwordValue := emailAndPassword[1]
	status, email := store.GetEmail(ctx, emailValue)

	if status != enums.Ok || email != nil {
		return
	}

	store.CreateEmailConfirmationCode(ctx, emailValue, environment.TestConfirmationCode, time.Minute)
	_, user := store.CreateUser(ctx, passwordValue, emailValue, "")
	_, role := store.GetAdminRole(ctx)
	store.CreateUserRole(ctx, user.Id, role.Id)
}

func InitialAdminRole(environment *config.Environment, store stores.IStore) {
	ctx := context.Background()
	_, role := store.GetAdminRole(ctx)
	if role != nil {
		_, _ = store.CreateRole(ctx, environment.AdminRole)
	}
}

func main() {

	// Initialization

	environment := config.InitEnvironment()
	tracer, tracerCloser := config.InitTracing(environment)
	config.InitSentry(environment)
	pool := config.InitPool(tracer, environment)
	redis := config.InitRedis(environment)
	inMemoryCache := config.InitInMemoryCache()
	producer := config.InitNSQProducer(environment)
	passwordProcessor := passwordProcessors.InitPasswordProcessor()
	postgresRepo := postgresRepository.InitPostgresRepository(pool)
	redisRepo := redisRepository.InitRedisRepository(redis)
	inMemoryRepo := inMemoryRepository.InitInMemoryRepository(inMemoryCache)
	store := stores.InitStore(pool, redis, inMemoryCache, environment, postgresRepo, redisRepo, inMemoryRepo)
	jwtAuthenticationBackend := backends.InitJWTAuthenticationBackend(store)
	basicAuthenticationBackend := backends.InitBasicAuthenticationBackend(store, passwordProcessor)
	authenticationController := auth.InitAuthController(map[string]backends.IAuthenticationBackend{
		"Basic":  basicAuthenticationBackend,
		"Bearer": jwtAuthenticationBackend,
	}, environment)
	dispatcher := eventDispatchers.InitNSQEventDispatcher(producer, environment)
	controller := controllers.InitController(store, passwordProcessor, dispatcher, environment)
	API := api2.InitAPI(controller, authenticationController, environment)
	InitialAdminRole(environment, store)
	InitialAdmin(environment, store)

	authentication := middlewares.AuthenticationMiddleware(authenticationController)
	isLocalRequest := middlewares.IsLocalRequestMiddleware(environment.LocalNetworkNamespace)

	// Init Routing

	router := mux.NewRouter().StrictSlash(false)

	CreateUserV1 := http.HandlerFunc(API.CreateUserV1)
	GetUsersV1 := authentication(http.HandlerFunc(API.GetUsersV1))
	GetUserV1 := authentication(http.HandlerFunc(API.GetUserV1))
	DeleteUserV1 := authentication(http.HandlerFunc(API.DeleteUserV1))

	CreatePasswordV1 := authentication(http.HandlerFunc(API.CreatePasswordV1))

	CreateEmailV1 := authentication(http.HandlerFunc(API.CreateEmailV1))
	CreateEmailConfirmationV1 := http.HandlerFunc(API.CreateEmailConfirmationV1)

	CreateRoleV1 := authentication(http.HandlerFunc(API.CreateRoleV1))
	GetRolesV1 := authentication(http.HandlerFunc(API.GetRolesV1))
	GetRoleV1 := authentication(http.HandlerFunc(API.GetRoleV1))

	CreateUserRoleV1 := authentication(http.HandlerFunc(API.CreateUserRoleV1))
	GetUserRolesV1 := authentication(http.HandlerFunc(API.GetUserRolesV1))
	DeleteUserRoleV1 := authentication(http.HandlerFunc(API.DeleteUserRoleV1))

	CreatePhoneConfirmationV1 := http.HandlerFunc(API.CreatePhoneConfirmationV1)
	CreatePhoneV1 := authentication(http.HandlerFunc(API.CreatePhoneV1))

	CreateSessionV1 := authentication(http.HandlerFunc(API.CreateSessionV1))

	GetSecretV1 := isLocalRequest(http.HandlerFunc(API.GetSecretV1))

	GetUserViewV1 := authentication(http.HandlerFunc(API.GetUserViewV1))
	GetUsersViewV1 := authentication(http.HandlerFunc(API.GetUsersViewV1))

	uuidRE := "[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}"

	router.Handle("/api/v1/users", CreateUserV1).Methods(http.MethodPost)
	router.Handle("/api/v1/users", GetUsersV1).Methods(http.MethodGet)
	router.Handle(fmt.Sprintf("/api/v1/users/{id:%s}", uuidRE), GetUserV1).Methods(http.MethodGet)
	router.Handle(fmt.Sprintf("/api/v1/users/{id:%s}", uuidRE), DeleteUserV1).Methods(http.MethodDelete)

	router.Handle("/api/v1/passwords", CreatePasswordV1).Methods(http.MethodPost)

	router.Handle("/api/v1/emails", CreateEmailV1).Methods(http.MethodPost)
	router.Handle("/api/v1/emailConfirmations", CreateEmailConfirmationV1).Methods(http.MethodPost)

	router.Handle("/api/v1/roles", CreateRoleV1).Methods(http.MethodPost)
	router.Handle("/api/v1/roles", GetRolesV1).Methods(http.MethodGet)
	router.Handle(fmt.Sprintf("/api/v1/roles/{id:%s}", uuidRE), GetRoleV1).Methods(http.MethodGet)

	router.Handle("/api/v1/userRoles", CreateUserRoleV1).Methods(http.MethodPost)
	router.Handle("/api/v1/userRoles", GetUserRolesV1).Methods(http.MethodGet)
	router.Handle(fmt.Sprintf("/api/v1/userRoles/{id:%s}", uuidRE), DeleteUserRoleV1).Methods(http.MethodDelete)

	router.Handle("/api/v1/phoneConfirmations", CreatePhoneConfirmationV1).Methods(http.MethodPost)
	router.Handle("/api/v1/phones", CreatePhoneV1).Methods(http.MethodPost)

	router.Handle("/api/v1/sessions", CreateSessionV1).Methods(http.MethodPost)

	router.Handle(fmt.Sprintf("/api/v1/secrets/{id:%s}", uuidRE), GetSecretV1).Methods(http.MethodGet)

	router.Handle("/views/v1/users", GetUserViewV1).Methods(http.MethodGet)
	router.Handle(fmt.Sprintf("/views/v1/users/{id:%s}", uuidRE), GetUsersViewV1).Methods(http.MethodGet)

	// Middleware

	router.Use(middlewares.TracerMiddleware(tracer))
	router.Use(sentryHttp.New(sentryHttp.Options{}).Handle)
	router.Use(middlewares.ContentTypeMiddleware)

	// Finish

	http.Handle("/", router)

	defer tracerCloser.Close()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
