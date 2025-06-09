package order

import (
	"github.com/go-redis/redis/v8"
)

type orderCacheRepository struct {
	cache *redis.Client
}

func NewOrderCacheRepository(cache *redis.Client) *orderCacheRepository {
	return &orderCacheRepository{cache: cache}
}
