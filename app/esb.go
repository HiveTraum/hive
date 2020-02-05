package app

import (
	"auth/config"
	"auth/controllers"
	"auth/events"
	"auth/infrastructure"
)

func InitESB(store infrastructure.StoreInterface) *controllers.ESB {
	env := config.InitEnv()
	dispatcher := events.EventDispatcher{Url: env.EsbUrl}
	return &controllers.ESB{Store: store, Dispatcher: &dispatcher}
}
