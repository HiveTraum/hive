package main

import (
	"context"
	"fmt"
	sentryHttp "github.com/getsentry/sentry-go/http"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose"
	"github.com/urfave/cli/v2"
	api2 "hive/api"
	"hive/auth"
	"hive/auth/backends"
	"hive/config"
	"hive/controllers"
	"hive/enums"
	"hive/eventDispatchers"
	"hive/middlewares"
	"hive/passwordProcessors"
	"hive/repositories/inMemoryRepository"
	"hive/repositories/postgresRepository"
	"hive/repositories/redisRepository"
	"hive/stores"
	"log"
	"net/http"
	"os"
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
	if role == nil {
		_, _ = store.CreateRole(ctx, environment.AdminRole)
	}
}

func server() {

	// Initialization

	environment := config.InitEnvironment()
	tracer, tracerCloser := config.InitTracing(environment)
	config.InitSentry(environment)
	pool := config.InitPool(tracer, environment)
	redis := config.InitRedis(environment)
	inMemoryCache := config.InitInMemoryCache()
	producer := config.InitNSQProducer(environment)
	passwordProcessor := passwordProcessors.InitPasswordProcessor()
	postgresRepo := postgresRepository.InitPostgresRepository(pool, environment)
	redisRepo := redisRepository.InitRedisRepository(redis)
	inMemoryRepo := inMemoryRepository.InitInMemoryRepository(inMemoryCache)
	store := stores.InitStore(pool, redis, inMemoryCache, environment, postgresRepo, redisRepo, inMemoryRepo)
	jwtAuthenticationBackend := backends.InitJWTAuthenticationBackend(store, environment)
	basicAuthenticationBackend := backends.InitBasicAuthenticationBackend(store, passwordProcessor, environment)
	authenticationController := auth.InitAuthController(map[string]backends.IAuthenticationBackend{
		"Basic":  basicAuthenticationBackend,
		"Bearer": jwtAuthenticationBackend,
	}, environment)
	dispatcher := eventDispatchers.InitNSQEventDispatcher(producer, environment)
	controller := controllers.InitController(store, passwordProcessor, dispatcher, environment, jwtAuthenticationBackend.EncodeAccessToken)
	API := api2.InitAPI(controller, authenticationController, environment)
	InitialAdminRole(environment, store)
	InitialAdmin(environment, store)

	authentication := middlewares.AuthenticationMiddleware(authenticationController)
	isLocalRequest := middlewares.IsLocalRequestMiddleware(environment.LocalNetworkNamespace)

	// Init Routing

	router := mux.NewRouter().StrictSlash(false)

	CreateUserV1 := http.HandlerFunc(API.CreateUserV1)
	GetUsersV1 := authentication(http.HandlerFunc(API.GetUsersV1), true)
	GetUserV1 := authentication(http.HandlerFunc(API.GetUserV1), true)
	DeleteUserV1 := authentication(http.HandlerFunc(API.DeleteUserV1), true)

	CreatePasswordV1 := authentication(http.HandlerFunc(API.CreatePasswordV1), true)

	CreateEmailV1 := authentication(http.HandlerFunc(API.CreateEmailV1), true)
	CreateEmailConfirmationV1 := http.HandlerFunc(API.CreateEmailConfirmationV1)

	CreateRoleV1 := authentication(http.HandlerFunc(API.CreateRoleV1), true)
	GetRolesV1 := authentication(http.HandlerFunc(API.GetRolesV1), true)
	GetRoleV1 := authentication(http.HandlerFunc(API.GetRoleV1), true)

	CreateUserRoleV1 := authentication(http.HandlerFunc(API.CreateUserRoleV1), true)
	GetUserRolesV1 := authentication(http.HandlerFunc(API.GetUserRolesV1), true)
	DeleteUserRoleV1 := authentication(http.HandlerFunc(API.DeleteUserRoleV1), true)

	CreatePhoneConfirmationV1 := http.HandlerFunc(API.CreatePhoneConfirmationV1)
	CreatePhoneV1 := authentication(http.HandlerFunc(API.CreatePhoneV1), true)

	CreateSessionV1 := authentication(http.HandlerFunc(API.CreateSessionV1), false)

	GetSecretV1 := isLocalRequest(http.HandlerFunc(API.GetSecretV1))

	GetUserViewV1 := authentication(http.HandlerFunc(API.GetUserViewV1), true)
	GetUsersViewV1 := authentication(http.HandlerFunc(API.GetUsersViewV1), true)

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

func migrate() error {
	environment := config.InitEnvironment()
	dbConfig, err := pgx.ParseConfig(environment.DatabaseURI)
	if err != nil {
		return err
	}
	db := stdlib.OpenDB(*dbConfig)
	err = goose.Up(db, "./migrations")
	if err != nil {
		return err
	}
	return db.Close()
}

func main() {
	app := &cli.App{
		Name:  "hive",
		Usage: "Start hive server",
		Action: func(c *cli.Context) error {
			server()
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:  "migrate",
				Usage: "Apply migrations for current database",
				Action: func(c *cli.Context) error {
					migrate()
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
