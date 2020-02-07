package app

import (
	"auth/config"
	"auth/controllers"
	"auth/infrastructure"
	"auth/stores"
	"github.com/opentracing/opentracing-go"
)

type App struct {
	Store           infrastructure.StoreInterface
	ESB             infrastructure.ESBInterface
	LoginController infrastructure.LoginControllerInterface
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
	pool := config.InitPool(tracer)
	redis := config.InitRedis()
	inMemoryCache := config.InitInMemoryCache()
	store := &stores.DatabaseStore{Db: pool, Cache: redis, InMemoryCache: inMemoryCache}
	loginController := &controllers.LoginController{Store: store}
	esb := InitESB(store)
	InitialAdminRole(store)
	InitialAdmin(store)
	return &App{ESB: esb, Store: store, LoginController: loginController}
}
