package config

import (
	"github.com/go-redis/redis/v7"
)

func InitRedis() *redis.Client {

	env := InitEnv()

	client := redis.NewClient(&redis.Options{
		Addr: env.RedisUrl,
	})

	_, err := client.Ping().Result()

	if err != nil {
		panic(err)
	}

	return client
}
