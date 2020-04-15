package app

import (
	"auth/backends"
	"auth/config"
	"auth/controllers"
	"auth/infrastructure"
	"auth/processors"
	"auth/stores"
	"github.com/opentracing/opentracing-go"
)

type App struct {
	Store             infrastructure.StoreInterface
	ESB               infrastructure.ESBInterface
	LoginController   infrastructure.LoginControllerInterface
	PasswordProcessor infrastructure.PasswordProcessorInterface
}

func (app *App) GetPasswordProcessor() infrastructure.PasswordProcessorInterface {
	return app.PasswordProcessor
}

func (app *App) GetStore() infrastructure.StoreInterface {
	return app.Store
}

func (app *App) GetESB() infrastructure.ESBInterface {
	return app.ESB
}

func (app *App) GetLoginController() infrastructure.LoginControllerInterface {
	return app.LoginController
}

func InitApp(tracer opentracing.Tracer) *App {
	config.InitSentry()
	env := config.GetEnvironment()
	pool := config.InitPool(tracer)
	redis := config.InitRedis()
	inMemoryCache := config.InitInMemoryCache()
	store := &stores.DatabaseStore{Db: pool, Cache: redis, InMemoryCache: inMemoryCache}
	authBackends := map[string]infrastructure.AuthenticationBackend{
		"Bearer": backends.JWTAuthenticationBackend{Store: store},
		"Basic":  backends.BasicAuthenticationBackend{Store: store},
	}
	loginController := &controllers.LoginController{Backends: authBackends, RequestContextUserKey: env.RequestContextUserKey}
	passwordProcessor := &processors.PasswordProcessor{}
	esb := InitESB(store)
	InitialAdminRole(store)
	InitialAdmin(store)
	return &App{ESB: esb, Store: store, LoginController: loginController, PasswordProcessor: passwordProcessor}
}
