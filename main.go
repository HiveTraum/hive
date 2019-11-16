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

	phoneConfirmationsAPI := middlewares.ResponseControllerMiddleware(api.PhoneConfirmationsAPIV1(application))
	phonesAPI := middlewares.ResponseControllerMiddleware(api.PhonesAPIV1(application))

	emailConfirmationsAPI := middlewares.ResponseControllerMiddleware(api.EmailConfirmationsAPIV1(application))
	emailsAPI := middlewares.ResponseControllerMiddleware(api.EmailsAPIV1(application))

	// Middleware

	// Methods Middleware
	CR := middlewares.IsMethodAllowedMiddleware([]string{http.MethodGet, http.MethodPost})
	R := middlewares.IsMethodAllowedMiddleware([]string{http.MethodGet})
	C := middlewares.IsMethodAllowedMiddleware([]string{http.MethodPost})

	handlers := []handler{
		{pattern: "/views/v1/users", h: R(usersView),},
		{pattern: "/views/v1/users/{id:[0-9]+}", h: R(userView),},

		{pattern: "/api/v1/users", h: CR(usersAPI),},
		{pattern: "/api/v1/users/{id:[0-9]+}", h: R(userAPI),},

		{pattern: "/api/v1/passwords", h: C(passwordsAPI),},

		{pattern: "/api/v1/roles", h: CR(rolesAPI),},
		{pattern: "/api/v1/roles/{id:[0-9]+}", h: R(roleAPI),},

		{pattern: "/api/v1/phoneConfirmations", h: C(phoneConfirmationsAPI),},
		{pattern: "/api/v1/phones", h: C(phonesAPI),},

		{pattern: "/api/v1/emailConfirmations", h: C(emailConfirmationsAPI),},
		{pattern: "/api/v1/emails", h: C(emailsAPI),},
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
		router.HandleFunc(h.pattern, h.h)
	}

	// Tracing routes

	_ = router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		route.Handler(gorilla.Middleware(tracer, route.GetHandler()))
		return nil
	})

	// Finish

	http.Handle("/", router)

	log.Fatal(http.ListenAndServe(":8080", nil))

	defer closer.Close()
}
