package envs

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type IMDBItem interface {
	GetIMDBType() IMDBType
	GetDBKey() string
	GetRedisClient() *redis.Client
	GetDB() any
	CloseDB(ctx context.Context) error
	GetConfig() any
}

type imdbItem struct {
	dbKey  string
	dbType IMDBType
	redis  *redis.Client
	config any
}

func newIMDBItem(dbKey string, dbType IMDBType, redis *redis.Client, config any) IMDBItem {
	return &imdbItem{
		dbKey:  dbKey,
		dbType: dbType,
		redis:  redis,
		config: config,
	}
}

func (i *imdbItem) GetIMDBType() IMDBType {
	return i.dbType
}

func (i *imdbItem) GetDBKey() string {
	return i.dbKey
}

func (i *imdbItem) GetRedisClient() *redis.Client {
	return i.redis
}

func (i *imdbItem) GetDB() any {
	return i.config
}

func (i *imdbItem) CloseDB(ctx context.Context) error {
	return i.redis.Close()
}

func (i *imdbItem) GetConfig() any {
	return i.config
}
