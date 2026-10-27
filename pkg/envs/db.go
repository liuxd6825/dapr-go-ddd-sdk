package envs

import (
	"strings"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/boolutils"
)

type DbsEnv struct {
	items    map[string]DBItem
	mongo    map[string]*Mongo
	mysql    map[string]*MySql
	doris    map[string]*Doris
	postgres map[string]*Postgres
	elastic  map[string]*Elastic
}

func newDbsEnv() *DbsEnv {
	return &DbsEnv{
		items:    map[string]DBItem{},
		mongo:    map[string]*Mongo{},
		mysql:    map[string]*MySql{},
		doris:    map[string]*Doris{},
		postgres: map[string]*Postgres{},
		elastic:  map[string]*Elastic{},
	}
}

func initDB(env *Env) {
	errList := errors.NewErrors()
	for dbKey, info := range env.DBs {
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
		case DBTypeDoris.String():
			db, err := NewDorisWithMap(info)
			if err != nil {
				errList.AddError(errors.New("dbKey %s error: %s", dbKey, err.Error()))
				continue
			}
			if err := AddDBDoris(env, db); err != nil {
				errList.AddError(err)
			}
			env.dbs.doris[db.DbKey] = db
			break
		case DBTypePostgres.String():
			db, err := NewPostgresWithMap(info)
			if err != nil {
				errList.AddError(errors.New("dbKey %s error: %s", dbKey, err.Error()))
				continue
			}
			if err := AddDBPostgres(env, db); err != nil {
				errList.AddError(err)
			}
			env.dbs.postgres[db.DbKey] = db
			break
		case DBTypeMySQL.String():
			db, err := NewMySQLWithMap(info)
			if err != nil {
				errList.AddError(errors.New("dbKey %s error: %s", dbKey, err.Error()))
				continue
			}
			if err := AddDBMySql(env, db); err != nil {
				errList.AddError(err)
			}
			env.dbs.mysql[db.DbKey] = db
			break
		case DBTypeMongoDB.String():
			db, err := NewMongoWithMap(info)
			if err != nil {
				errList.AddError(errors.New("dbKey %s error: %s", dbKey, err.Error()))
				continue
			}
			if err := AddDBMongo(env, db); err != nil {
				errList.AddError(err)
			}

			env.dbs.mongo[db.DbKey] = db
			break
		case DBTypeElastic.String():
			db, err := NewElasticWithMap(info)
			if err != nil {
				errList.AddError(errors.New("dbKey %s error: %s", dbKey, err.Error()))
				continue
			}
			if err := AddDBElastic(env, db); err != nil {
				errList.AddError(err)
			}
			env.dbs.elastic[db.DbKey] = db
			break
		default:
			errList.AddError(errors.New("dbType is invalid, ", dbKey))
		}

	}

	if errList.HasError() {
		panicError(errList)
	}
}

func (d DBType) String() string {
	return string(d)
}
