package envs

import (
	"context"

	"github.com/apache/hugegraph-toolchain/hugegraph-client-go"
	"github.com/elastic/go-elasticsearch/v9"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/mongodb"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/gp"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/redis/go-redis/v9"
	nebula "github.com/vesoft-inc/nebula-go/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type DBItem interface {
	GetDBType() DBType
	GetDBKey() string
	GetNeo4jClient() neo4j.DriverWithContext
	GetMongoClient() *mongo.Client
	GetMongoDB() *mongodb.MongoDB
	GetGormDBClient() *gorm.DB
	GetElasticClient() *elasticsearch.Client
	GetHugeClient() *hugegraph.CommonClient
	GetRedisClient() *redis.Client
	GetNebulaClient() *nebula.ConnectionPool
	GetDB() any
	CloseDB(ctx context.Context) error
	GetConfig() any
}

type EventPublish interface {
	GetEventPublish() bool
}

type dbItem struct {
	dbKey      string
	dbType     DBType
	redis      *redis.Client
	neo4j      neo4j.DriverWithContext
	mongo      *mongodb.MongoDB
	gormDb     *gorm.DB
	elastic    *elasticsearch.Client
	hugeClient *hugegraph.CommonClient
	nebulaPool *nebula.ConnectionPool
	config     any
}

type DBType string

const (
	DBTypePostgres DBType = "postgres"
	DBTypeMySQL    DBType = "mysql"
	DBTypeSqlite   DBType = "sqlite"
	DBTypeNeo4j    DBType = "neo4j"
	DBTypeMongoDB  DBType = "mongodb"
	DBTypeMsSQL    DBType = "mssql"
	DBTypeOracle   DBType = "oracle"
	DBTypeElastic  DBType = "elastic"
	DBTypeHuge     DBType = "hugegraph"
	DBTypeDoris    DBType = "doris"
	DBTypeNebula   DBType = "nebula"
	DBTypeRedis    DBType = "redis"
)

func (d *dbItem) GetDB() any {
	switch d.dbType {
	case DBTypePostgres:
		return d.gormDb
	case DBTypeMySQL:
		return d.gormDb
	case DBTypeSqlite:
		return d.gormDb
	case DBTypeNeo4j:
		return d.neo4j
	case DBTypeMongoDB:
		return d.mongo
	case DBTypeMsSQL:
		return d.gormDb
	case DBTypeOracle:
		return d.gormDb
	case DBTypeElastic:
		return d.elastic
	case DBTypeHuge:
		return d.hugeClient
	case DBTypeDoris:
		return d.gormDb
	case DBTypeNebula:
		return d.nebulaPool
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

func (d *dbItem) GetNeo4jClient() neo4j.DriverWithContext {
	return d.neo4j
}

func (d *dbItem) GetMongoClient() *mongo.Client {
	return d.mongo.GetClient()
}

func (d *dbItem) GetMongoDB() *mongodb.MongoDB {
	return d.mongo
}

func (d *dbItem) GetGormDBClient() *gorm.DB {
	return d.gormDb
}

func (d *dbItem) GetElasticClient() *elasticsearch.Client {
	return d.elastic
}

func (d *dbItem) GetHugeClient() *hugegraph.CommonClient {
	return d.hugeClient
}

func (d *dbItem) GetRedisClient() *redis.Client {
	return d.redis
}

func (d *dbItem) GetNebulaClient() *nebula.ConnectionPool {
	return d.nebulaPool
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
