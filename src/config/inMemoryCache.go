package config

import (
	"github.com/patrickmn/go-cache"
	"time"
)

func InitInMemoryCache() *cache.Cache {
	return cache.New(time.Minute*5, time.Minute*10)
}
