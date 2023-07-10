package cache

import (
	"time"
	"trade-robot-bd/libs/cache"
)

func CacheData(key string, data interface{}, duration time.Duration) error {
	return cache.Redis().Set(key, data, duration).Err()
}

func GetData(key string) ([]byte, error) {
	return cache.Redis().Get(key).Bytes()
}
