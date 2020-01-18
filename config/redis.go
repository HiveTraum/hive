package config

import (
	"github.com/go-redis/redis/v7"
	"sync"
)

var client *redis.Client
var onceRedis sync.Once

func InitRedis() *redis.Client {

	env := InitEnv()
	onceRedis.Do(func() {
		client = redis.NewClient(&redis.Options{
			Addr: env.RedisUrl,
		})

		_, err := client.Ping().Result()

		if err != nil {
			panic(err)
		}
	})

	return client
}
