package sql

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/fsm"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/db_pkg/db"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/schema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/gp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/idutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/randomutils"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
	"time"
)

type Human struct {
	Id         string    `gorm:"primaryKey"`
	Name       string    `gorm:"name"`
	Age        int       `gorm:"age"`
	Analyse    string    `gorm:"analyse"`
	Birthday   time.Time `gorm:"birthday"`
	PeopleType []string  `gorm:"people_type;type:text[]"`
}

func Test_DB(t *testing.T) {
	database, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		t.Fatalf("数据库连接失败: %v", err)
		return
	}
	human := &Human{
		Id:         idutils.NewId(),
		PeopleType: []string{"Human", "People", "PeopleType"},
	}
	table := database.Model(human)
	table.AutoMigrate(human)
	res := table.Create(human)

	assert.Equal(t, 1, res.RowsAffected)
}

func Test_Dao(t *testing.T) {
	database, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		t.Fatalf("数据库连接失败: %v", err)
		return
	}

	humanName := randomutils.NameCN()
	human := map[string]any{
		"id":         idutils.NewId(),
		"tenantId":   "test",
		"analyse":    "",
		"birthday":   time.Now(),
		"peopleType": []string{"1111"},
		"name":       humanName,
		"age":        1,
		"tags":       []string{"tag1", "tag2"},
	}
	t.Log(human)

	humanSchema, err := schema.NewSchemaWithJson("human.json", HumanSchema)
	if err != nil {
		t.Error(err)
		return
	}

	daoCfg := &db.DaoConfig{
		Database:   database,
		DbKey:      "sql",
		Schema:     humanSchema.GetJsonSchema(),
		Env:        NewEnvConfig(),
		IsPubEvent: false,
	}

	dao := NewDao(daoCfg)
	ctx, err := restapp.NewTestContext(context.Background())
	if err != nil {
		t.Error(err)
		return
	}

	/*
		t.Run("dao.AutoMigrate", func(t *testing.T) {
			gp.Try(func() error {
				dao.Table().AutoMigrate(ctx)
				return nil
			}).Catch(func(e error) {
				t.Error(err)
			})
		})
	*/

	t.Run("dao.Create", func(t *testing.T) {
		gp.Try(func() error {
			dao.Create(ctx, human)
			return nil
		}).Catch(func(err error) {
			t.Error(err)
		})
	})

	t.Run("dao.Update", func(t *testing.T) {
		gp.Try(func() error {
			humanName = humanName + "2"
			human["name"] = humanName
			dao.Update(ctx, human)
			return nil
		}).Catch(func(err error) {
			t.Error(err)
		})
	})

	t.Run("dao.FindByRSQL", func(t *testing.T) {
		list := dao.FindByRSQL(ctx, fmt.Sprintf("name=='%s'", humanName))
		t.Log("list:", list)
	})

	t.Run("dao.FindAll", func(t *testing.T) {
		res := dao.FindAll(ctx)
		if res.Error != nil {
			t.Error(res.Error)
		} else {
			t.Log("list:", res.GetData())
		}
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
