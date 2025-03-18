package restapp

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository/ddd_mongodb"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type DbType string

const (
	DbType_Postgres DbType = "postgres"
	DbType_MySQL    DbType = "mysql"
	DbType_Sqlite   DbType = "sqlite"
	DbType_Neo4j    DbType = "neo4j"
	DbType_Redis    DbType = "redis"
	DbType_MongoDB  DbType = "mongodb"
	DbType_MsSQL    DbType = "mssql"
	DbType_Oracle   DbType = "oracle"
)

type DBItem interface {
	GetDBType() DbType
	GetDBKey() string
	GetRedis() *redis.Client
	GetNeo4j() neo4j.DriverWithContext
	GetMongo() *ddd_mongodb.MongoDB
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
	dbType DbType
	redis  *redis.Client
	neo4j  neo4j.DriverWithContext
	mongo  *ddd_mongodb.MongoDB
	gormDb *gorm.DB
	config any
}

var _dbs map[string]DBItem

func init() {
	_dbs = make(map[string]DBItem)
}

func (d *dbItem) GetDB() any {
	switch d.dbType {
	case DbType_Postgres:
		return d.gormDb
	case DbType_MySQL:
		return d.gormDb
	case DbType_Sqlite:
		return d.gormDb
	case DbType_Neo4j:
		return d.neo4j
	case DbType_Redis:
		return d.redis
	case DbType_MongoDB:
		return d.mongo
	case DbType_MsSQL:
		return d.gormDb
	case DbType_Oracle:
		return d.gormDb
	default:
		panic("db type not supported")
	}
}

func (d *dbItem) GetDBType() DbType {
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

func (d *dbItem) GetMongo() *ddd_mongodb.MongoDB {
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

func addDb(dbItem DBItem) {
	dbKey := dbItem.GetDBKey()
	_, ok := _dbs[dbKey]
	if ok {
		panic(errors.New("db \"%s\" already exists", dbKey))
	}
	_dbs[dbKey] = dbItem
}

func CloseAllDb(ctx context.Context) error {
	for _, d := range _dbs {
		_ = d.CloseDB(ctx)
	}
	return nil
}

func GetDb(dbKey string) DBItem {
	item := _dbs[dbKey]
	return item
}
