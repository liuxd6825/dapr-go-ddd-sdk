package test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/iris-contrib/httpexpect/v2"
	"github.com/kataras/iris/v12/httptest"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"testing"
)

var (
	EnvName                 = "dev_lxd"
	TenantId                = "test"
	DefaultOperationTimeout = 5
)

type Option struct {
	fileName   *string
	envName    *string
	eventTypes []restapp.RegisterEventType
}

func NewOption(options ...*Option) *Option {
	envName := EnvName
	o := &Option{envName: &envName}
	for _, item := range options {
		if item.fileName != nil {
			o.fileName = item.fileName
		}
		if item.eventTypes != nil {
			o.eventTypes = item.eventTypes
		}
		if item.envName != nil {
			o.envName = item.envName
		}
	}
	return o
}

func InitQuery(options ...*Option) error {
	o := NewOption(options...)
	if o.fileName == nil {
		o.SetFileName("query-config.yaml")
	}
	if o.envName == nil {
		o.SetEnvName(EnvName)
	}
	ctx := context.Background()
	return Init(ctx, o.FileName(), o.EnvName(), o.EventTypes())
}

func InitCommand(options ...*Option) error {
	o := NewOption(options...)
	if o.fileName == nil {
		o.SetFileName("cmd-config.yaml")
	}
	if o.envName == nil {
		o.SetEnvName(EnvName)
	}
	ctx := context.Background()
	return Init(ctx, o.FileName(), o.EnvName(), o.EventTypes())
}

// InitQuery
// @Description: 初始化QueryService测试环境
func Init(ctx context.Context, fileName string, envName string, eventTypes []restapp.RegisterEventType) error {
	if fileName == "" {
		return errors.New("fileName is null")
	}
	if envName == "" {
		return errors.New("envName is null")
	}
	// 加载配置文件
	config, err := restapp.NewConfigByFile("${search}/config/" + fileName)
	if err != nil {
		return err
	}

	// 读取指定环境下的配置信息
	env, err := config.GetEnvConfig(envName)
	if err != nil {
		return err
	}

	return restapp.InitApplication(ctx, env, eventTypes, nil)
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

/*
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
*/

/*
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
*/

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

func (o *Option) EnvName() string {
	if o.envName == nil {
		return EnvName
	}
	return *o.envName
}

func (o *Option) FileName() string {
	if o.fileName == nil {
		return ""
	}
	return *o.fileName
}

func (o *Option) EventTypes() []restapp.RegisterEventType {
	if o.eventTypes == nil {
		return nil
	}
	return o.eventTypes
}

func (o *Option) SetEnvName(val string) *Option {
	o.envName = &val
	return o
}

func (o *Option) SetFileName(val string) *Option {
	o.fileName = &val
	return o
}

func (o *Option) SetEventTypes(val []restapp.RegisterEventType) *Option {
	o.eventTypes = val
	return o
}
