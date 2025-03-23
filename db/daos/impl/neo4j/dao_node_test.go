package neo4j

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/core/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/daos/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/db/rsql"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/schema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types/times"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/gp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/idutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/randomutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/xtest"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_NodeDao(t *testing.T) {
	humanName := randomutils.NameCN()
	humanSchema, err := schema.NewSchemaWithJson("human.json", xtest.HumanSchema)
	if err != nil {
		t.Error(err)
		return
	}

	daoCfg := &idao.DaoConfig{
		Database:   driver,
		DbKey:      "neo4j",
		DBSchema:   dbschema.NewDBSchemaWithJsonSchema(humanSchema.GetJsonSchema()),
		Env:        xtest.NewEnvConfig(),
		IsPubEvent: false,
	}

	dao := NewDao[map[string]any](daoCfg)
	ctx, err := restapp.NewTestContext(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	dao.Table().Drop(ctx)
	dao.Table().AutoMigrate(ctx)
	newCount := int64(10)
	var list []map[string]any

	for i := int64(0); i < newCount; i++ {
		entity := map[string]any{
			"id":         randomutils.NewId(),
			"name":       humanName,
			"analyse":    "",
			"age":        randomutils.IntMax(100),
			"birthday":   times.NowTime(),
			"peopleType": []string{"1111"},
			"tags":       []string{"tag1", "tag2"},
			"caseId":     "test",
			"graphId":    "neo4j-test",
		}
		list = append(list, entity)
	}

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
				v["name"] = humanName + "UpdateMany"
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
		"caseId":     "test",
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

	/*
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

	*/

	t.Run("dao.FindById", func(t *testing.T) {
		gp.Try(func() error {
			entity := dao.FindById(ctx, id)
			assert.NotNil(t, entity)
			if entity != nil {
				if dataId, ok := entity["id"].(string); ok {
					assert.Equal(t, id, dataId)
				}
				t.Log("findById:", entity)
				if _, ok := entity["peopleType"]; !ok {
					t.Error("peopleType not exist")
				}
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
			assert.NoError(t, res.Error)
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
			res := dao.UpdateByRSQL(ctx, builder.Build(), human)
			assert.Equal(t, int64(1), res.RowsAffected)
			return nil
		}).Catch(func(err error) {
			t.Error(err)
		})
	})

	return

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

	t.Run("dao.DeleteByRSQL", func(t *testing.T) {
		gp.Try(func() error {
			res := dao.DeleteByRSQL(ctx, fmt.Sprintf("creatorName=='%s'", "test"))
			t.Log("DeleteByRSQL count:", res.RowsAffected)
			assert.Equal(t, newCount, res.RowsAffected)
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
