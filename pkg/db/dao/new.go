package dao

import (
	"fmt"
	"runtime"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/impl/huge"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/impl/mongodb"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/impl/neo4j"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/impl/sql"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/types"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/reflectutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/stringutils"
)

type DaoConfig struct {
	DBKey string `json:"dbKey"` // 可选 DBKey 与 DB
	DB    any    `json:"db"`    // 可选 DBKey 与 DB
	// IsPubEvent         *bool           `json:"isPubEvent"`         // 可选
	AggField           string          `json:"aggField"`           // 可选
	AggType            string          `json:"aggType"`            // 可选
	DBSchema           *store.DBSchema `json:"dbSchema"`           // 可选
	TableName          string          `json:"tableName"`          // 可选
	GraphType          idao.GraphType  `json:"graphType"`          // 可选
	GraphLabels        []string        `json:"graphLabels"`        // 可选
	IsCancelModified   bool            `json:"isCancelModified"`   // 可选
	IsCancelSoftDelete bool            `json:"isCancelSoftDelete"` // 可选
	Env                *envs.Env       `json:"-"`                  // 可选
	IsCache            bool            `json:"isCache"`
}

var cache = types.NewCMap[any]()

var isFalse bool = false
var isTrue bool = true

func IsFalse() *bool {
	return &isFalse
}

func IsTrue() *bool {
	return &isTrue
}

func NewDao[T any](newCfg *DaoConfig) idao.Dao[T] {
	if newCfg == nil {
		panic("new dao must have a non-nil pointer newCfg parameter")
	}
	if newCfg.DBKey == "" {
		panic("new dao must have a non-empty db key at newCfg.DBKey parameter")
	}

	dbKey := newCfg.DBKey
	tableName := newCfg.TableName
	env := newCfg.Env
	if env == nil {
		env = envs.GetEnv()
	}
	dbKey = GetDbKey(env, dbKey)

	if tableName == "" && newCfg.DBSchema != nil {
		tableName = newCfg.DBSchema.TableName
	}
	var daoKey string
	if newCfg.IsCache {
		// dao缓存 取得daoKey
		className := reflectutils.GetClassName[T]()
		daoKey = getDaoKey(dbKey, tableName, className)
		if v, ok := cache.Get(daoKey); v != nil && ok {
			return v.(idao.Dao[T])
		}
	}

	// 是struct类型
	if newCfg.DBSchema == nil && reflectutils.IsStruct[T]() {
		entity := reflectutils.NewInstance[T]()
		if tableName == "" {
			tableName = reflectutils.GetStructName(entity)
			tableName = stringutils.SnakeString(tableName)
		}
		dbSchema := dbschema.NewDBSchemaWithStruct(tableName, entity, tableName)
		newCfg.DBSchema = dbSchema
	}

	if newCfg.DBSchema == nil {
		panic("db.NewDao() args schema is required")
	}

	item := getDBItem(env, dbKey)
	if item == nil && newCfg.DB != nil {
		newCfg.DB = newCfg.DB
	} else if item == nil {
		panic(fmt.Sprintf("pkg.db.dao.NewDao() error: No dbKey exists \"%s\" in Env \"%s\", \"%s\"", dbKey, env.Name, getCallerInfo(2)))
	}

	/*
		isPubEvent := false
		if newCfg.IsPubEvent != nil {
			isPubEvent = *newCfg.IsPubEvent
		} else {
			isPubEvent = env.GetEnv().App.IsPubEvent
		}*/

	daoCfg := &idao.DaoConfig{
		DBKey: dbKey,
		DB:    newCfg.DB,
		//IsPubEvent:         isPubEvent,
		AggField:           newCfg.AggField,
		AggType:            newCfg.AggType,
		Env:                env,
		DBSchema:           newCfg.DBSchema,
		GraphType:          newCfg.GraphType,
		GraphLabels:        newCfg.GraphLabels,
		IsCancelModified:   newCfg.IsCancelModified,
		IsCancelSoftDelete: newCfg.IsCancelSoftDelete,
	}

	/*if isPubEvent {
		daoCfg.OutboxDao = NewOutboxDao(newCfg.DBKey)
	}*/

	var dao idao.Dao[T]

	if item != nil {
		switch item.GetDBType() {
		case envs.DBTypeMongoDB:
			dao = mongodb.NewDao[T](daoCfg)
			break
		case envs.DBTypeSqlite, envs.DBTypeOracle, envs.DBTypePostgres, envs.DBTypeMySQL, envs.DBTypeDoris, envs.DBTypeMsSQL:
			dao = sql.NewDao[T](daoCfg)
			break
		case envs.DBTypeNeo4j:
			dao = neo4j.NewDao[T](daoCfg)
			break
		case envs.DBTypeHuge:
			dao = huge.NewDao[T](daoCfg)
			break
		default:
			panic(errors.New(fmt.Sprintf("dao.NewDao[T any]() not supported dbKey %s database %s, tableName:%s", dbKey, item.GetDBType().String(), newCfg.TableName)))
		}
	}
	if newCfg.IsCache {
		cache.Add(daoKey, any(dao))
	}
	return dao
}

func getCallerInfo(skip int) string {
	// skip = 1 表示获取上一层调用者的信息（0表示当前函数）
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}
	return fmt.Sprintf("%s:%d", file, line)
}

func GetDbKey(envInst *envs.Env, dbKey string) string {
	return envInst.GetDBKeyValue(dbKey)
}

func getDaoKey(dbKey, tableName string, className string) string {
	return fmt.Sprintf("%s.%s.%s", dbKey, tableName, className)
}

func getDBItem(env *envs.Env, dbKey string) envs.DBItem {
	item := env.GetDB(dbKey)
	return item
}

func NewConfig(dbKey string, tableName string, entity any) *DaoConfig {
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, entity, tableName)
	newCfg := &DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	return newCfg
}
