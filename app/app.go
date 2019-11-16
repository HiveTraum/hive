package app

import (
	"auth/config"
	"auth/controllers"
	"auth/events"
	"auth/infrastructure"
	"auth/stores"
	"github.com/opentracing/opentracing-go"
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

func InitESB(store *stores.DatabaseStore) *controllers.ESB {
	env := config.InitEnv()
	dispatcher := events.EventDispatcher{Url: env.EsbUrl}
	return &controllers.ESB{Store: store, Dispatcher: &dispatcher}
}

func InitApp(tracer opentracing.Tracer) *App {
	config.InitSentry()
	pool := config.InitPool(tracer)
	redis := config.InitRedis()
	store := stores.DatabaseStore{Db: pool, Cache: redis,}
	esb := InitESB(&store)
	return &App{ESB: esb, Store: &store}
}
