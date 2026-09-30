package hugedao

import (
	"context"
	"os"
	"testing"

	"github.com/apache/hugegraph-toolchain/hugegraph-client-go"
	"github.com/elastic/go-elasticsearch/v9"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/mongodb"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"gorm.io/gorm"
)

const testGraphName = "dev"
const testDbKey = "huge"

var testClient *hugegraph.CommonClient

func TestMain(m *testing.M) {
	if v := os.Getenv("HUGEGRAPH_SKIP_INIT"); v == "" {
		client, err := hugegraph.NewCommonClient(hugegraph.Config{
			Host:     envOr("HUGEGRAPH_HOST", "192.168.120.200"),
			Port:     envOrInt("HUGEGRAPH_PORT", 18080),
			Graph:    testGraphName,
			Username: "admin",
			Password: "admin",
		})
		if err == nil {
			if res, verr := client.Version(); verr == nil && res != nil {
				if res.Body != nil {
					_ = res.Body.Close()
				}
				testClient = client
			}
		}
	}
	os.Exit(m.Run())
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envOrInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	n := 0
	for _, r := range v {
		if r < '0' || r > '9' {
			return def
		}
		n = n*10 + int(r-'0')
	}
	return n
}

// testDBItem 实现 env.DBItem,只暴露 hugedao client,其它字段为空
type testDBItem struct {
	dbKey  string
	client *hugegraph.CommonClient
}

func (d *testDBItem) GetDBType() env.DBType             { return env.DBType_Huge }
func (d *testDBItem) GetDBKey() string                  { return d.dbKey }
func (d *testDBItem) GetNeo4j() neo4j.DriverWithContext { return nil }
func (d *testDBItem) GetMongo() *mongodb.MongoDB        { return nil }
func (d *testDBItem) GetGormDB() *gorm.DB               { return nil }
func (d *testDBItem) GetElastic() *elasticsearch.Client { return nil }
func (d *testDBItem) GetHuge() *hugegraph.CommonClient  { return d.client }
func (d *testDBItem) GetDB() any                        { return d.client }
func (d *testDBItem) CloseDB(_ context.Context) error   { return nil }
func (d *testDBItem) GetConfig() any                    { return nil }

func registerTestEnv(client *hugegraph.CommonClient) {
	e := env.GetEnv()
	if e == nil {
		e = env.NewEnv()
		env.SetEnv(e)
	}
	if e.GetDB(testDbKey) == nil {
		e.AddDB(&testDBItem{dbKey: testDbKey, client: client})
	}
}
