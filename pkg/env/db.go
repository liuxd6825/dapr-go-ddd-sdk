package env

import (
	"context"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/mongodb"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/gp"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type DBItem interface {
	GetDBType() DBType
	GetDBKey() string
	GetNeo4j() neo4j.DriverWithContext
	GetMongo() *mongodb.MongoDB
	GetGormDB() *gorm.DB
	GetElastic() *elasticsearch.Client
	GetDB() any
	CloseDB(ctx context.Context) error
	GetConfig() any
}

type EventPublish interface {
	GetEventPublish() bool
}

type dbItem struct {
	dbKey   string
	dbType  DBType
	redis   *redis.Client
	neo4j   neo4j.DriverWithContext
	mongo   *mongodb.MongoDB
	gormDb  *gorm.DB
	elastic *elasticsearch.Client
	config  any
}

type DBType string

const DefaultDBKey string = "default"

const (
	DBType_Postgres DBType = "postgres"
	DBType_MySQL    DBType = "mysql"
	DBType_Sqlite   DBType = "sqlite"
	DBType_Neo4j    DBType = "neo4j"
	DBType_MongoDB  DBType = "mongodb"
	DBType_MsSQL    DBType = "mssql"
	DBType_Oracle   DBType = "oracle"
	DBType_Elastic  DBType = "elastic"
)

func (d DBType) String() string {
	return string(d)
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
	case DBType_MongoDB:
		return d.mongo
	case DBType_MsSQL:
		return d.gormDb
	case DBType_Oracle:
		return d.gormDb
	case DBType_Elastic:
		return d.elastic
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

func (d *dbItem) GetNeo4j() neo4j.DriverWithContext {
	return d.neo4j
}

func (d *dbItem) GetMongo() *mongodb.MongoDB {
	return d.mongo
}

func (d *dbItem) GetGormDB() *gorm.DB {
	return d.gormDb
}

func (d *dbItem) GetElastic() *elasticsearch.Client {
	return d.elastic
}

func (d *dbItem) CloseDB(ctx context.Context) error {
	return gp.Try(func() error {
		if d.mongo != nil {
			return d.mongo.Close(ctx)
		}
		if d.neo4j != nil {
			return d.neo4j.Close(ctx)
		}
		if d.gormDb != nil {
			return d.gormClose(d.gormDb)
		}
		if d.elastic != nil {
			return d.elastic.Close(ctx)
		}
		return nil
	}).Error
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
