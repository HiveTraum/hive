package config

import (
	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog/log"
	"time"
)

func InitInMemoryCache() *cache.Cache {
	c := cache.New(time.Minute*5, time.Minute*10)
	log.Log().Msg("In memory cache successfully initiated")
	return c
}
