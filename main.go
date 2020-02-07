package main

import (
	"auth/api"
	"auth/app"
	"auth/middlewares"
	"auth/views"
	sentryHttp "github.com/getsentry/sentry-go/http"
	"github.com/gorilla/mux"
	"github.com/opentracing-contrib/go-gorilla/gorilla"
	"github.com/opentracing/opentracing-go"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
	"log"
	"net/http"
)

type handler struct {
	pattern string
	h       http.HandlerFunc
	methods []string
}

func main() {

	// Tracing

	tracer, closer, err := jaegerConfig.Configuration{
		ServiceName: "auth",
		RPCMetrics:  true,
	}.NewTracer()

	if err != nil {
		panic(err)
	}

	opentracing.SetGlobalTracer(tracer)

	// App

	application := app.InitApp(tracer)

	// Handlers

	usersView := middlewares.ResponseControllerMiddleware(views.UsersViewV1(application))
	userView := middlewares.ResponseControllerMiddleware(views.UserViewV1(application))

	usersAPI := middlewares.ResponseControllerMiddleware(api.UsersAPIV1(application))
	userAPI := middlewares.ResponseControllerMiddleware(api.UserAPIV1(application))

	passwordsAPI := middlewares.ResponseControllerMiddleware(api.PasswordsAPIV1(application))

	rolesAPI := middlewares.ResponseControllerMiddleware(api.RolesAPIV1(application))
	roleAPI := middlewares.ResponseControllerMiddleware(api.RoleAPIV1(application))

	userRolesAPI := middlewares.ResponseControllerMiddleware(api.UserRolesAPIV1(application))
	userRoleAPI := middlewares.ResponseControllerMiddleware(api.UserRoleAPIV1(application))

	phoneConfirmationsAPI := middlewares.ResponseControllerMiddleware(api.PhoneConfirmationsAPIV1(application))
	phonesAPI := middlewares.ResponseControllerMiddleware(api.PhonesAPIV1(application))

	emailConfirmationsAPI := middlewares.ResponseControllerMiddleware(api.EmailConfirmationsAPIV1(application))
	emailsAPI := middlewares.ResponseControllerMiddleware(api.EmailsAPIV1(application))

	sessionsAPI := middlewares.ResponseControllerMiddleware(api.SessionsAPIV1(application))

	secretsAPI := middlewares.ResponseControllerMiddleware(api.SecretAPIV1(application))

	// Middleware

	// Methods Middleware
	CR := []string{http.MethodGet, http.MethodPost, http.MethodOptions}
	RD := []string{http.MethodGet, http.MethodDelete, http.MethodOptions}
	R := []string{http.MethodGet, http.MethodOptions}
	C := []string{http.MethodPost, http.MethodOptions}
	D := []string{http.MethodDelete, http.MethodOptions}

	handlers := []handler{
		{pattern: "/views/v1/users", h: usersView, methods: R},
		{pattern: "/views/v1/users/{id:[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}}", h: userView, methods: RD},

		{pattern: "/api/v1/users", h: usersAPI, methods: CR},
		{pattern: "/api/v1/users/{id:[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}}", h: userAPI, methods: R},

		{pattern: "/api/v1/passwords", h: passwordsAPI, methods: C},

		{pattern: "/api/v1/roles", h: rolesAPI, methods: CR},
		{pattern: "/api/v1/roles/{id:[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}}", h: roleAPI, methods: R},

		{pattern: "/api/v1/userRoles", h: userRolesAPI, methods: CR},
		{pattern: "/api/v1/userRoles/{id:[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}}", h: userRoleAPI, methods: D},

		{pattern: "/api/v1/phoneConfirmations", h: phoneConfirmationsAPI, methods: C},
		{pattern: "/api/v1/phones", h: phonesAPI, methods: C},

		{pattern: "/api/v1/emailConfirmations", h: emailConfirmationsAPI, methods: C},
		{pattern: "/api/v1/emails", h: emailsAPI, methods: C},

		{pattern: "/api/v1/sessions", h: sessionsAPI, methods: C},

		{pattern: "/api/v1/secrets/{id:[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}}", h: secretsAPI, methods: R},
	}

	// Content Type Middleware

	for i, h := range handlers {
		h.h = middlewares.ContentTypeMiddleware(h.h)
		handlers[i] = h
	}

	// Sentry

	sentryHandler := sentryHttp.New(sentryHttp.Options{})

	for i, h := range handlers {
		h.h = sentryHandler.HandleFunc(h.h)
		handlers[i] = h
	}

	// Init Routing

	router := mux.NewRouter().StrictSlash(false)

	for _, h := range handlers {
		router.HandleFunc(h.pattern, h.h).Methods(h.methods...)
	}

	// Tracing routes

	_ = router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		route.Handler(gorilla.Middleware(tracer, route.GetHandler()))
		return nil
	})

	// Finish

	http.Handle("/", router)

	defer closer.Close()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
