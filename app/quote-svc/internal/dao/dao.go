package dao

import (
	"github.com/go-redis/redis"
	"trade-robot-bd/libs/cache"
)

type Dao struct {
	RedisCli *redis.Client
}

func New() *Dao {
	return &Dao{
		RedisCli: cache.Redis(),
	}
}
