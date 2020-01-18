package stores

import (
	"github.com/go-redis/redis/v7"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/patrickmn/go-cache"
)

type DatabaseStore struct {
	Db            *pgxpool.Pool
	Cache         *redis.Client
	InMemoryCache *cache.Cache
}
