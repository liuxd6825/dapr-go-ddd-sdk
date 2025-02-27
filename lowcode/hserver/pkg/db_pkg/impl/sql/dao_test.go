package sql

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
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
	dao.Table().Drop(ctx)
	dao.Table().AutoMigrate(ctx)
	newCount := int64(10)
	var list []map[string]any
	t.Run("dao.CreateMany", func(t *testing.T) {
		gp.Try(func() error {
			for i := int64(0); i < newCount; i++ {
				entity := map[string]any{
					"id":         randomutils.NewId(),
					"name":       humanName,
					"analyse":    "",
					"age":        randomutils.IntMax(100),
					"birthday":   randomutils.Date(),
					"peopleType": []string{"1111"},
					"tags":       []string{"tag1", "tag2"},
				}
				list = append(list, entity)
			}
			dao.CreateMany(ctx, list)
			return nil
		}).Catch(func(err error) {
			t.Error(err)
		})
	})

	t.Run("dao.UpdateMany", func(t *testing.T) {
		gp.Try(func() error {
			for _, v := range list {
				v["remark"] = "remark," + randomutils.String(10)
			}
			dao.UpdateMany(ctx, list)
			return nil
		}).Catch(func(err error) {
			t.Error(err)
		})
	})

	id := idutils.NewId()
	human := map[string]any{
		"id":         id,
		"tenantId":   "test",
		"analyse":    "",
		"birthday":   time.Now(),
		"peopleType": []string{"1111"},
		"name":       humanName,
		"age":        1,
		"tags":       []string{"tag1", "tag2"},
	}

	t.Run("dao.Create", func(t *testing.T) {
		gp.Try(func() error {
			dao.Create(ctx, human)
			return nil
		}).Catch(func(err error) {
			t.Error(err)
		})
	})

	t.Run("dao.FindById", func(t *testing.T) {
		gp.Try(func() error {
			e := dao.FindById(ctx, id)
			t.Log("findById:", e)
			return nil
		}).Catch(func(err error) {
			t.Error(err)
		})
	})

	t.Run("dao.Update", func(t *testing.T) {
		gp.Try(func() error {
			human["name"] = humanName + "2"
			count := dao.Update(ctx, human)
			assert.Equal(t, int64(1), count)
			return nil
		}).Catch(func(err error) {
			t.Error(err)
		})
	})

	t.Run("dao.UpdateByRSQL", func(t *testing.T) {
		gp.Try(func() error {
			humanName = humanName + "3"
			human["name"] = humanName
			count := dao.UpdateByRSQL(ctx, fmt.Sprintf("id=='%s'", id), human)
			assert.Equal(t, int64(1), count)
			return nil
		}).Catch(func(err error) {
			t.Error(err)
		})
	})

	t.Run("dao.DeleteById", func(t *testing.T) {
		gp.Try(func() error {
			dao.DeleteById(ctx, id)
			t.Log("deleteById:", id)
			return nil
		}).Catch(func(err error) {
			t.Error(err)
		})
	})

	t.Run("dao.FindByRSQL", func(t *testing.T) {
		findList := dao.FindByRSQL(ctx, fmt.Sprintf("creatorName=='%s'", "test"))
		t.Log("list:", findList)
		assert.Equal(t, newCount, int64(len(findList)))
	})

	t.Run("dao.FindAll", func(t *testing.T) {
		res := dao.FindAll(ctx)
		if res.Error != nil {
			t.Error(res.Error)
		} else {
			t.Log("list:", res.GetData())
		}
	})

	t.Run("dao.CountByRSQL", func(t *testing.T) {
		gp.Try(func() error {
			res := dao.CountByRSQL(ctx, fmt.Sprintf("creatorName=='%s'", "test"))
			assert.Equal(t, newCount, res)
			t.Log("count:", res)
			return nil
		}).Catch(func(err error) {
			t.Error(err)
		})
	})

	t.Run("dao.Sum", func(t *testing.T) {
		gp.Try(func() error {
			var vals []*ddd_repository.ValueCol
			vals = append(vals, &ddd_repository.ValueCol{
				AggFunc: "sum",
				Field:   "age",
			})
			res := dao.SumByRSQL(ctx, "", vals)
			t.Log("count:", res)
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
