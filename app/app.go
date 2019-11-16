package app

import (
	"auth/config"
	"auth/controllers"
	"auth/infrastructure"
	"auth/stores"
)

type App struct {
	Store infrastructure.StoreInterface
	ESB   infrastructure.ESBInterface
}

func (app *App) GetStore() infrastructure.StoreInterface {
	return app.Store
}

func (app *App) GetESB() infrastructure.ESBInterface {
	return app.ESB
}

func InitApp() *App {
	config.InitSentry()
	pool := config.InitPool()
	redis := config.InitRedis()
	store := stores.DatabaseStore{Db: pool, Cache: redis,}
	esb := &controllers.ESB{Store: &store}
	return &App{ESB: esb, Store: &store}
}
