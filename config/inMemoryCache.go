package config

import (
	"github.com/patrickmn/go-cache"
	"sync"
	"time"
)

var inMemoryCache *cache.Cache
var onceInMemoryCache sync.Once

func InitInMemoryCache() *cache.Cache {

	onceInMemoryCache.Do(func() {
		inMemoryCache = cache.New(time.Minute*5, time.Minute*10)
	})

	return inMemoryCache
}
