package envs

import (
	"strings"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/boolutils"
)

type IMDBType string

const (
	IMDBTypeRedis IMDBType = "redis"
)

type IMDBEnv struct {
	items map[string]IMDBItem
	redis map[string]*Redis
}

func newIMDBEnv() *IMDBEnv {
	return &IMDBEnv{
		items: make(map[string]IMDBItem),
		redis: map[string]*Redis{},
	}
}

func initIMDB(env *Env) {
	errList := errors.NewErrors()
	for dbKey, info := range env.IMDB {
		dbType, isHas := info["dbType"]
		info["dbKey"] = strings.ToLower(dbKey)
		if !isHas {
			errList.AddError(errors.New("%s dbType is required", dbKey))
		}
		enable := boolutils.GetBool(info["enable"], true)
		if !enable {
			continue
		}

		switch dbType {
		case IMDBTypeRedis.String():
			db, err := NewRedisWithMap(info)
			if err != nil {
				errList.AddError(errors.New("dbKey %s error: %s", dbKey, err.Error()))
				continue
			}
			if err := AddIMDBRedis(env, db); err != nil {
				errList.AddError(err)
			}
			env.imdb.redis[db.DbKey] = db
			break
		default:
			errList.AddError(errors.New("dbType is invalid, ", dbKey))
		}
	}
	if errList.HasError() {
		panicError(errList)
	}
}

func (t IMDBType) String() string {
	return string(t)
}
