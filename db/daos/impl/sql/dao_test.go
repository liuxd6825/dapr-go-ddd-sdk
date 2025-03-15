package sql

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/daos/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/daos/impl/tests"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/schema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/rsql"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types/times"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/gp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/idutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/randomutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/reflectutils"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
	"time"
)

type Human struct {
	Id         string
	TenantId   string
	Name       string
	Age        int
	Analyse    string
	Birthday   *time.Time
	PeopleType string
	Tags       string
}

func Test_DaoStruct(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Error(err)
		}
	}()
	ctx, err := restapp.NewTestContext(context.Background())
	if err != nil {
		t.Fatal(err)
		return
	}
	dao := newDaoByStruct[*Human](ctx, "human_struct")
	humanName := randomutils.NameCN()
	newCount := int64(10)
	list := newStructList(newCount, humanName)
	dao.CreateMany(ctx, list)
}

func Test_Dao(t *testing.T) {
	ctx, err := restapp.NewTestContext(context.Background())
	if err != nil {
		t.Fatal(err)
		return
	}
	dao := newDao[map[string]any](ctx, "human")
	humanName := randomutils.NameCN()
	newCount := int64(10)
	list := newMapList(newCount, humanName)

	t.Run("dao.DeleteByRSQL", func(t *testing.T) {
		gp.Try(func() error {
			list := newMapList(newCount, "0000")
			res := dao.CreateMany(ctx, list)
			assert.Equal(t, newCount, res.RowsAffected)

			delRes := dao.DeleteByRSQL(ctx, "name=='0000'")
			t.Log("DeleteByRSQL count:", delRes.RowsAffected)
			assert.Equal(t, newCount, delRes.RowsAffected)

			return nil
		}).Catch(func(err error) {
			t.Error(err)
		})
	})

	t.Run("dao.CreateMany", func(t *testing.T) {
		gp.Try(func() error {
			res := dao.CreateMany(ctx, list)
			assert.Equal(t, newCount, res.RowsAffected)
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
			res := dao.UpdateMany(ctx, list)
			assert.Equal(t, newCount, res.RowsAffected)
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
			res := dao.Create(ctx, human)
			assert.Equal(t, int64(1), res.RowsAffected)
			return nil
		}).Catch(func(err error) {
			t.Error(err)
		})
	})

	t.Run("dao.Update", func(t *testing.T) {
		gp.Try(func() error {
			human["name"] = humanName + "2"
			human["birthday"] = times.NewDate()
			count := dao.Update(ctx, human)
			assert.Equal(t, int64(1), count.RowsAffected)

			human["birthday"] = times.NewTime()
			count = dao.Update(ctx, human)
			assert.Equal(t, int64(1), count.RowsAffected)

			return nil
		}).Catch(func(err error) {
			t.Error(err)
		})
	})

	t.Run("dao.FindById", func(t *testing.T) {
		gp.Try(func() error {
			e := dao.FindById(ctx, id)
			assert.NotNil(t, e)
			if e != nil {
				if dataId, ok := e["id"].(string); ok {
					assert.Equal(t, id, dataId)
				}
			}
			t.Log("findById:", e)

			if _, ok := e["peopleType"]; !ok {
				t.Error("peopleType not exist")
			}

			return nil
		}).Catch(func(err error) {
			t.Error(err)
		})
	})

	t.Run("dao.FindPaging", func(t *testing.T) {
		gp.Try(func() error {
			paging := ddd_repository.NewFindPagingQueryRequest()
			paging.PageSize = 2
			paging.IsTotalRows = true
			paging.Filter = fmt.Sprintf("creatorName=='%s'", "test")
			res := dao.FindPaging(ctx, paging)
			assert.NotNil(t, res)
			assert.Equal(t, int64(11), res.TotalRows)
			assert.Equal(t, 2, len(res.Data))
			t.Log("totalRows=", res.TotalRows, " count=", len(res.Data))
			return nil
		}).Catch(func(err error) {
			t.Error(err)
		})
	})

	t.Run("dao.UpdateByRSQL", func(t *testing.T) {
		gp.Try(func() error {
			humanName = humanName + "3"
			human["name"] = humanName
			builder := rsql.NewBuilder().Eq("id", id)
			count := dao.UpdateByRSQL(ctx, builder.Build(), human)
			assert.Equal(t, int64(1), count.RowsAffected)
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
		builder := rsql.NewBuilder().Eq("creatorName", "test")
		findList := dao.FindByRSQL(ctx, builder.Build())
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

	t.Run("dao.SumByRSQL", func(t *testing.T) {
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

	t.Run("dao.DeleteAll", func(t *testing.T) {
		gp.Try(func() error {
			res1 := dao.CreateMany(ctx, list)
			assert.Equal(t, newCount, res1.RowsAffected)

			res2 := dao.DeleteAll(ctx)
			t.Log("DeleteAll count:", res2.RowsAffected)
			assert.Equal(t, newCount, res2.RowsAffected)
			return nil
		}).Catch(func(err error) {
			t.Error(err)
		})
	})

	t.Run("dao.DeleteByIds", func(t *testing.T) {
		gp.Try(func() error {
			res1 := dao.CreateMany(ctx, list)
			assert.Equal(t, newCount, res1.RowsAffected)

			var ids []string
			for _, e := range list {
				if v, ok := e["id"].(string); ok {
					ids = append(ids, v)
				}
			}

			res2 := dao.DeleteByIds(ctx, ids)
			t.Log("DeleteByIds count:", res2.RowsAffected)
			assert.Equal(t, newCount, res2.RowsAffected)
			return nil
		}).Catch(func(err error) {
			t.Error(err)
		})
	})

}

func TestDao_Sum(t *testing.T) {
	ctx, err := restapp.NewTestContext(context.Background())
	if err != nil {
		t.Fatal(err)
		return
	}
	dao := newDao[map[string]any](ctx, "human_sum")
	if dao == nil {
		return
	}

	if count := dao.CountByRSQL(ctx, ""); count == 0 {
		var list []map[string]any
		newCount := int64(10)

		for i := int64(0); i < newCount; i++ {
			entity := map[string]any{
				"id":         randomutils.NewId(),
				"name":       "sum",
				"analyse":    "",
				"age":        randomutils.IntMax(100),
				"birthday":   randomutils.Date(),
				"peopleType": []string{"1111"},
				"tags":       []string{"tag1", "tag2"},
			}
			list = append(list, entity)
		}
		res := dao.CreateMany(ctx, list)
		assert.Equal(t, newCount, res.RowsAffected)
	}

	t.Run("sum", func(t *testing.T) {
		valueCols := make([]*ddd_repository.ValueCol, 0)
		valueCols = append(valueCols, &ddd_repository.ValueCol{
			AggFunc: "sum", Field: "age",
		})
		qry := ddd_repository.NewFindPagingQueryRequest()
		qry.SetTenantId("test")
		qry.SetPageSize(2)
		qry.SetValueCols(valueCols)

		sumAny := dao.SumByQuery(ctx, qry)
		fmt.Println("data:", sumAny)
	})

	t.Run("group", func(t *testing.T) {
		valueCols := make([]*ddd_repository.ValueCol, 0)
		valueCols = append(valueCols, &ddd_repository.ValueCol{
			AggFunc: "sum", Field: "age",
		})

		groupCols := make([]*ddd_repository.GroupCol, 0)
		groupCols = append(groupCols, &ddd_repository.GroupCol{
			Field: "gender", DataType: "string",
		})

		groupKeys := make([]any, 0)

		qry := ddd_repository.NewFindPagingQueryRequest()
		qry.SetTenantId("test")
		qry.SetPageSize(2)
		qry.SetValueCols(valueCols)
		qry.SetGroupCols(groupCols)
		qry.SetGroupKeys(groupKeys)
		qry.SetIsTotalRows(true)

		findRes := dao.FindPaging(ctx, qry)
		fmt.Println("data:", findRes)

	})
}

func TestDao_Find(t *testing.T) {
	ctx, err := restapp.NewTestContext(context.Background())
	if err != nil {
		t.Fatal(err)
		return
	}
	dao := newDao[map[string]any](ctx, "human_sum")

	dao.DeleteAll(ctx)
	newCount := int64(10)
	newList := newMapList(newCount, "000")
	res := dao.CreateMany(ctx, newList)
	assert.Equal(t, newCount, res.RowsAffected)

	t.Run("birthday>'2000-01-01'", func(t *testing.T) {
		list := dao.FindByRSQL(ctx, "birthday>'2000-01-01'")
		t.Log(len(list))
	})

}

func newStructList(count int64, humanName string) []*Human {
	list := make([]*Human, count)
	for i := int64(0); i < count; i++ {
		entity := &Human{
			Id:         randomutils.NewId(),
			Name:       humanName,
			Analyse:    "",
			Age:        randomutils.IntMax(100),
			Birthday:   randomutils.PDate(),
			PeopleType: randomutils.String(10),
			Tags:       randomutils.String(10),
		}
		list[i] = entity
	}
	return list
}

func newMapList(count int64, humanName string) []map[string]any {
	list := make([]map[string]any, count)
	for i := int64(0); i < count; i++ {
		entity := map[string]any{
			"id":         randomutils.NewId(),
			"name":       humanName,
			"analyse":    "",
			"age":        randomutils.IntMax(100),
			"birthday":   randomutils.Date(),
			"peopleType": []string{"1111"},
			"tags":       []string{"tag1", "tag2"},
		}
		list[i] = entity
	}
	return list
}

func newDaoByStruct[T any](ctx context.Context, tableName string) idao.Dao[T] {
	database, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("数据库连接失败: %v", err))
	}
	data := reflectutils.NewInstance[T]()
	dbSch := dbschema.NewSchemaWithStruct(tableName, data, tableName)
	daoCfg := &idao.DaoConfig{
		Database:   database,
		DbKey:      "sql",
		Schema:     dbSch,
		Env:        tests.NewEnvConfig(),
		IsPubEvent: false,
	}

	dao := NewDao[T](daoCfg)
	dao.Table().Drop(ctx)
	dao.Table().AutoMigrate(ctx)
	return dao
}

func newDao[T any](ctx context.Context, tableName string) idao.Dao[T] {
	database, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("数据库连接失败: %v", err))
	}

	//humanName := randomutils.NameCN()
	humanSchema, err := schema.NewSchemaWithJson("human.json", tests.HumanSchema)
	if err != nil {
		panic(err)
	}

	dbSch := dbschema.NewSchemaWithJsonSchema(humanSchema.GetJsonSchema())
	dbSch.TableName = tableName

	daoCfg := &idao.DaoConfig{
		Database:   database,
		DbKey:      "sql",
		Schema:     dbSch,
		Env:        tests.NewEnvConfig(),
		IsPubEvent: false,
	}

	dao := NewDao[T](daoCfg)
	dao.Table().Drop(ctx)
	dao.Table().AutoMigrate(ctx)
	return dao
}
