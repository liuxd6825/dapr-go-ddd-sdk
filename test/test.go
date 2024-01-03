package test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/iris-contrib/httpexpect/v2"
	"github.com/kataras/iris/v12/httptest"
	"github.com/liuxd6825/dapr-go-ddd-sdk/applog"
	"github.com/liuxd6825/dapr-go-ddd-sdk/daprclient"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/dao/mongo_dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/dao/neo4j_dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository/ddd_mongodb"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/redis/go-redis/v9"
	"testing"
)

const (
	EnvName  = "dev_lxd"
	TenantId = "test"

	DefaultOperationTimeout = 5
)

// InitQuery
// @Description: 初始化QueryService测试环境
func InitQuery() error {
	// 加载配置文件
	config, err := restapp.NewConfigByFile("${search}/config/query-config.yaml")
	if err != nil {
		return err
	}

	// 读取指定环境下的配置信息
	env, err := config.GetEnvConfig(EnvName)
	if err != nil {
		return err
	}

	// 创建dapr客户端
	daprClient, err := daprclient.NewDaprDddClient(context.Background(), env.Dapr.GetHost(), env.Dapr.GetHttpPort(), env.Dapr.GetGrpcPort())
	if err != nil {
		return err
	}
	// 注册dapr客户端
	daprclient.SetDaprDddClient(daprClient)

	// 初始化ddd
	ddd.Init(env.App.AppId)
	applog.Init(daprClient, env.App.AppId, logs.DebugLevel)

	restapp.SetCurrentEnvConfig(env)
	// 初始化数据库
	if err := initMongo(env.Mongo["default"]); err != nil {
		return err
	}
	if err := initNeo4j(env.Neo4j["default"]); err != nil {
		return err
	}
	if err := initRedis(env.Redis["default"]); err != nil {
		return err
	}
	if err := initMinio(env); err != nil {
		return err
	}
	return nil
}

// GetResponseData
// @Description: 获取http响应数据
// @param t  *testing.T 测试对象
// @param resp *httpexpect.Response Http响应对象
// @param data 数据指针
// @return error 错误
func GetResponseData(t *testing.T, resp *httpexpect.Response, data interface{}) error {
	raw := resp.Raw()
	switch raw.StatusCode {
	case httptest.StatusOK:
		break
	case httptest.StatusNotFound:
		return errors.New(fmt.Sprintf("%v", raw.Status))
	default:
		return errors.New(fmt.Sprintf("%v, %v ", raw.Status, resp.Body().Raw()))
	}

	body := resp.Body()
	bytes := []byte(body.Raw())
	t.Log(body.Raw())
	if err := json.Unmarshal(bytes, data); err != nil {
		return err
	}
	return nil
}

func initMongo(dbConfig *restapp.MongoConfig) error {
	config := &ddd_mongodb.Config{
		Host:             dbConfig.Host,
		UserName:         dbConfig.UserName,
		Password:         dbConfig.Password,
		DatabaseName:     dbConfig.Database,
		OperationTimeout: DefaultOperationTimeout,
		ReplicaSet:       dbConfig.ReplicaSet,
	}
	db, err := ddd_mongodb.NewMongoDB(config, nil)
	if err != nil {
		return err
	}
	mongo_dao.SetDB(db)
	return nil
}

func initRedis(cfg *restapp.RedisConfig) error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Host,
		Password: cfg.Password, // no password set
		DB:       cfg.Database, // use default DB
	})
	restapp.SetRedis("default", rdb)
	return nil
}

func initMinio(cfg *restapp.EnvConfig) error {
	return restapp.InitMinioByEnvConfig(cfg)
}

func initNeo4j(dbConfig *restapp.Neo4jConfig) error {
	cfg := func(config *neo4j.Config) {
		config.MaxConnectionPoolSize = 10
	}
	neo4jUil := fmt.Sprintf("bolt://%v:%v", dbConfig.Host, dbConfig.Port)
	configures := []func(*neo4j.Config){cfg}
	driver, err := neo4j.NewDriverWithContext(neo4jUil, neo4j.BasicAuth(dbConfig.UserName, dbConfig.Password, ""), configures...)
	if err != nil {
		return err
	}
	neo4j_dao.SetDB(driver)
	return nil
}

// PrintJson 以JSON格式打印
func PrintJson(t *testing.T, title string, data interface{}) error {
	bs, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}
	t.Log(title)
	t.Log(string(bs))
	return nil
}
