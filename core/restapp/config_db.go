package restapp

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store/store_mongodb"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type DBType string

const (
	DBType_Postgres DBType = "postgres"
	DBType_MySQL    DBType = "mysql"
	DBType_Sqlite   DBType = "sqlite"
	DBType_Neo4j    DBType = "neo4j"
	DBType_Redis    DBType = "redis"
	DBType_MongoDB  DBType = "mongodb"
	DBType_MsSQL    DBType = "mssql"
	DBType_Oracle   DBType = "oracle"
)

type DBItem interface {
	GetDBType() DBType
	GetDBKey() string
	GetRedis() *redis.Client
	GetNeo4j() neo4j.DriverWithContext
	GetMongo() *store_mongodb.MongoDB
	GetGormDB() *gorm.DB
	GetDB() any
	CloseDB(ctx context.Context) error
	GetConfig() any
}

type EventPublish interface {
	GetEventPublish() bool
}

type dbItem struct {
	dbKey  string
	dbType DBType
	redis  *redis.Client
	neo4j  neo4j.DriverWithContext
	mongo  *store_mongodb.MongoDB
	gormDb *gorm.DB
	config any
}

var _dbs map[string]DBItem

func init() {
	_dbs = make(map[string]DBItem)
}

func (d *dbItem) GetDB() any {
	switch d.dbType {
	case DBType_Postgres:
		return d.gormDb
	case DBType_MySQL:
		return d.gormDb
	case DBType_Sqlite:
		return d.gormDb
	case DBType_Neo4j:
		return d.neo4j
	case DBType_Redis:
		return d.redis
	case DBType_MongoDB:
		return d.mongo
	case DBType_MsSQL:
		return d.gormDb
	case DBType_Oracle:
		return d.gormDb
	default:
		panic("db type not supported")
	}
}

func (d *dbItem) GetDBType() DBType {
	return d.dbType
}

func (d *dbItem) GetDBKey() string {
	return d.dbKey
}

func (d *dbItem) GetRedis() *redis.Client {
	return d.redis
}

func (d *dbItem) GetNeo4j() neo4j.DriverWithContext {
	return d.neo4j
}

func (d *dbItem) GetMongo() *store_mongodb.MongoDB {
	return d.mongo
}

func (d *dbItem) GetGormDB() *gorm.DB {
	return d.gormDb
}

func (d *dbItem) CloseDB(ctx context.Context) error {
	if d.mongo != nil {
		return d.mongo.Close(ctx)
	}
	if d.neo4j != nil {
		return d.neo4j.Close(ctx)
	}
	if d.redis != nil {
		return d.redis.Close()
	}
	if d.gormDb != nil {
		return d.gormClose(d.gormDb)
	}
	return nil
}

func (d *dbItem) gormClose(gormDb *gorm.DB) error {
	if gormDb != nil {
		db, err := gormDb.DB()
		if err != nil {
			return err
		}
		return db.Close()
	}
	return nil
}

func (d *dbItem) GetConfig() any {
	return d.config
}

func addDB(dbItem DBItem) {
	dbKey := dbItem.GetDBKey()
	_, ok := _dbs[dbKey]
	if ok {
		panic(errors.New("db \"%s\" already exists", dbKey))
	}
	_dbs[dbKey] = dbItem
}

func CloseAllDB(ctx context.Context) error {
	for _, d := range _dbs {
		_ = d.CloseDB(ctx)
	}
	return nil
}

func GetDB(dbKey string) DBItem {
	item := _dbs[dbKey]
	return item
}

func GetDBDefault() DBItem {
	for _, d := range _dbs {
		return d
	}
	return nil
}
