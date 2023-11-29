package localcache

import (
	"github.com/allegro/bigcache"
	"time"
)

var cache *bigcache.BigCache

func init() {
	cache, _ = bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
}

func SetCache(key string, data []byte) error {
	return cache.Set(key, data)
}

func GetCache(key string) ([]byte, error) {
	return cache.Get(key)
}

func GetString(key string) (*string, error) {
	b, err := cache.Get(key)
	var val string
	if err != nil {
		return &val, err
	}
	val = string(b)
	return &val, nil
}
