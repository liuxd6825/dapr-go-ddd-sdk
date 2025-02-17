package sql

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/fsm"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/db_pkg/db"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/schema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/gp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/idutils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
	"time"
)

type Human struct {
	Name       string    `gorm:"name"`
	Age        int       `gorm:"age"`
	Analyse    string    `gorm:"analyse"`
	Birthday   time.Time `gorm:"birthday"`
	PeopleType []string  `gorm:"people_type;type:json"`
}

func Test_Dao(t *testing.T) {
	database, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		t.Fatalf("数据库连接失败: %v", err)
		return
	}

	humanSchema, err := schema.NewSchemaWithJson("human.json", HumanSchema)
	if err != nil {
		t.Error(err)
		return
	}

	daoCfg := &db.DaoConfig{
		Database: database,
		DbKey:    "sql",
		Schema:   humanSchema,
		Env:      NewEnvConfig(),
	}

	dao := NewDao(daoCfg)
	ctx, err := restapp.NewTestContext(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	human := MapEntity{
		"id":         idutils.NewId(),
		"analyse":    "",
		"birthday":   "",
		"peopleType": []string{"1111"},
		"name":       "name",
		"age":        "age",
		"tags":       []string{"tag1", "tag2"},
	}

	t.Run("dao.AutoMigrate", func(t *testing.T) {
		gp.Try(func() error {
			dao.Table().AutoMigrate(ctx)
			return nil
		}).Catch(func(e error) {
			t.Error(err)
		})
	})

	t.Run("dao.Create", func(t *testing.T) {
		gp.Try(func() error {
			dao.Create(ctx, human)
			return nil
		}).Catch(func(err error) {
			t.Error(err)
		})
	})

}

type EnvConfig struct {
}

func (e EnvConfig) GetAppId() string {
	return "test"
}

func (e EnvConfig) GetAppName() string {
	return "app"
}

func (e EnvConfig) GetAppHttpHost() string {
	return "localhost"
}

func (e EnvConfig) GetAppHttpPort() int {
	return 0
}

func (e EnvConfig) GetDaprHost() string {
	return "localhost"
}

func (e EnvConfig) GetDaprHttpPort() int64 {
	return 0
}

func (e EnvConfig) GetDaprGrpcPort() int64 {
	return 0
}

func (e EnvConfig) GetFsManager() *fsm.Manager {
	return nil
}

func (e EnvConfig) GetHServerSrcPath() string {
	return ""
}

func (e EnvConfig) GetHServerEnable() bool {
	return false
}

func NewEnvConfig() *EnvConfig {
	return &EnvConfig{}
}
