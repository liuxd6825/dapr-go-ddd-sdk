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
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/randomutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/xtest"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_RelDao(t *testing.T) {
	humanName := randomutils.NameCN()

	relSchema, err := schema.NewSchemaWithJson("humanRel.json", xtest.HumanRelSchema)
	if err != nil {
		t.Error(err)
		return
	}

	relCfg := &idao.DaoConfig{
		Database:   driver,
		DbKey:      "neo4j",
		DBSchema:   dbschema.NewDBSchemaWithJsonSchema(relSchema.GetJsonSchema()),
		Env:        xtest.NewEnvConfig(),
		IsPubEvent: false,
	}

	relDao := NewDao[map[string]any](relCfg)
	ctx, err := restapp.NewTestContext(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	relDao.Table().Drop(ctx)
	relDao.Table().AutoMigrate(ctx)

	id := randomutils.NewId()
	nodeDao := newNodeDao(t)
	nodes := newNodes(humanName, id, t)
	nodeRes := nodeDao.CreateMany(ctx, nodes)
	assert.Equal(t, int64(len(nodes)), nodeRes.RowsAffected)

	rel := map[string]any{
		"id":         id,
		"startId":    "start" + id,
		"endId":      "end" + id,
		"relType":    "create",
		"tenantId":   "test",
		"analyse":    "",
		"birthday":   time.Now(),
		"peopleType": []string{"1111"},
		"name":       humanName,
		"age":        1,
		"tags":       []string{"tag1", "tag2"},
	}

	t.Run("rel.Create", func(t *testing.T) {
		gp.Try(func() error {
			res := relDao.Create(ctx, rel)
			assert.Equal(t, int64(1), res.RowsAffected)
			return nil
		}).Catch(func(err error) {
			t.Error(err)
		})
	})

	t.Run("rel.Update", func(t *testing.T) {
		gp.Try(func() error {
			rel["name"] = humanName + "-update"
			rel["birthday"] = times.NewDate()
			count := relDao.Update(ctx, rel)
			assert.Equal(t, int64(1), count.RowsAffected)
			return nil
		}).Catch(func(err error) {
			t.Error(err)
		})
	})

	t.Run("rel.FindById", func(t *testing.T) {
		gp.Try(func() error {
			entity := relDao.FindById(ctx, id)
			assert.NotNil(t, entity)
			return nil
		}).Catch(func(err error) {
			t.Error(err)
		})
	})

	t.Run("dao.UpdateByRSQL", func(t *testing.T) {
		gp.Try(func() error {
			humanName = humanName + "3"
			rel["name"] = humanName
			builder := rsql.NewBuilder().Eq("id", id)
			res := relDao.UpdateByRSQL(ctx, builder.Build(), rel)
			assert.Equal(t, int64(1), res.RowsAffected)
			return nil
		}).Catch(func(err error) {
			t.Error(err)
		})
	})

	t.Run("dao.DeleteById", func(t *testing.T) {
		gp.Try(func() error {
			relDao.DeleteById(ctx, id)
			t.Log("deleteById:", id)
			return nil
		}).Catch(func(err error) {
			t.Error(err)
		})
	})

}

func TestRelDao_Many(t *testing.T) {
	humanName := randomutils.NameCN()

	relSchema, err := schema.NewSchemaWithJson("humanRel.json", xtest.HumanRelSchema)
	if err != nil {
		t.Error(err)
		return
	}

	relCfg := &idao.DaoConfig{
		Database:   driver,
		DbKey:      "neo4j",
		DBSchema:   dbschema.NewDBSchemaWithJsonSchema(relSchema.GetJsonSchema()),
		Env:        xtest.NewEnvConfig(),
		IsPubEvent: false,
		DaoType:    "rel",
	}

	relDao := NewDao[map[string]any](relCfg, "rel0", "rel1")
	ctx, err := restapp.NewTestContext(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	relDao.Table().Drop(ctx)
	relDao.Table().AutoMigrate(ctx)
	newCount := int64(2)

	id := randomutils.NewId()
	nodeDao := newNodeDao(t)
	nodes := newNodes(humanName, id, t)

	var rels []map[string]any
	for i := int64(0); i < newCount; i++ {
		entity := map[string]any{
			"id":         randomutils.NewId(),
			"startId":    "start" + id,
			"endId":      "end" + id,
			"relType":    fmt.Sprintf("rel%d", i),
			"name":       humanName,
			"analyse":    "",
			"age":        randomutils.IntMax(100),
			"peopleType": []string{"1111"},
			"tags":       []string{"tag1", "tag2"},
			"caseId":     "test",
			"graphId":    "neo4j-test",
		}
		rels = append(rels, entity)
	}

	t.Run("rel.CreateMany", func(t *testing.T) {
		gp.Try(func() error {
			nodeRes := nodeDao.CreateMany(ctx, nodes)
			assert.Equal(t, int64(len(nodes)), nodeRes.RowsAffected)

			relRes := relDao.CreateMany(ctx, rels)
			assert.Equal(t, int64(len(rels)), relRes.RowsAffected)
			return nil
		}).Catch(func(err error) {
			t.Error(err)
		})
	})

	t.Run("rel.UpdateMany", func(t *testing.T) {
		gp.Try(func() error {
			for _, v := range rels {
				v["remark"] = "remark," + randomutils.String(10)
			}
			res := relDao.UpdateMany(ctx, rels)
			assert.Equal(t, newCount, res.RowsAffected)
			return nil
		}).Catch(func(err error) {
			t.Error(err)
		})
	})

	t.Run("rel.FindPaging", func(t *testing.T) {
		gp.Try(func() error {
			paging := ddd_repository.NewFindPagingQueryRequest()
			paging.PageSize = 2
			paging.IsTotalRows = true
			paging.Filter = "creatorName=='test'"
			res := relDao.FindPaging(ctx, paging)
			assert.NotNil(t, res)
			assert.NoError(t, res.Error)
			assert.Equal(t, 2, len(res.Data))
			t.Log("totalRows=", res.TotalRows, " count=", len(res.Data))
			return nil
		}).Catch(func(err error) {
			t.Error(err)
		})
	})

	t.Run("rel.FindByRSQL", func(t *testing.T) {
		rSql := rsql.NewBuilder().Eq("creatorName", "test").Build()
		findList := relDao.FindByRSQL(ctx, rSql)
		t.Log("list:", findList)
		//assert.Equal(t, newCount, int64(len(findList)))
	})

	t.Run("rel.FindAll", func(t *testing.T) {
		res := relDao.FindAll(ctx)
		if res.Error != nil {
			t.Error(res.Error)
		} else {
			t.Log("list:", res.GetData())
		}
	})

	t.Run("rel.CountByRSQL", func(t *testing.T) {
		gp.Try(func() error {
			res := relDao.CountByRSQL(ctx, fmt.Sprintf("creatorName=='%s'", "test"))
			t.Log("count:", res)
			return nil
		}).Catch(func(err error) {
			t.Error(err)
		})
	})

	t.Run("rel.Sum", func(t *testing.T) {
		gp.Try(func() error {
			var vals []*ddd_repository.ValueCol
			vals = append(vals, &ddd_repository.ValueCol{
				AggFunc: "sum",
				Field:   "age",
			})
			res := relDao.SumByRSQL(ctx, "", vals)
			t.Log("sum :", res)
			return nil
		}).Catch(func(err error) {
			t.Error(err)
		})
	})

	t.Run("rel.DeleteByRSQL", func(t *testing.T) {
		gp.Try(func() error {
			res := relDao.DeleteByRSQL(ctx, fmt.Sprintf("creatorName=='%s'", "test"))
			t.Log("DeleteByRSQL count:", res.RowsAffected)
			assert.Equal(t, newCount, res.RowsAffected)
			return nil
		}).Catch(func(err error) {
			t.Error(err)
		})
	})

	t.Run("rel.DeleteAll", func(t *testing.T) {
		gp.Try(func() error {
			res2 := relDao.DeleteAll(ctx)
			t.Log("DeleteAll count:", res2.RowsAffected)
			return nil
		}).Catch(func(err error) {
			t.Error(err)
		})
	})

	t.Run("rel.DeleteByIds", func(t *testing.T) {
		gp.Try(func() error {
			res1 := relDao.CreateMany(ctx, rels)
			assert.Equal(t, newCount, res1.RowsAffected)

			var ids []string
			for _, e := range rels {
				if v, ok := e["id"].(string); ok {
					ids = append(ids, v)
				}
			}

			res2 := relDao.DeleteByIds(ctx, ids)
			t.Log("DeleteByIds count:", res2.RowsAffected)
			assert.Equal(t, newCount, res2.RowsAffected)
			return nil
		}).Catch(func(err error) {
			t.Error(err)
		})
	})
}

func newNodeDao(t *testing.T) idao.Dao[map[string]any] {
	nodeSchema, err := schema.NewSchemaWithJson("human.json", xtest.HumanSchema)
	if err != nil {
		t.Error(err)
		return nil
	}

	nodeCfg := &idao.DaoConfig{
		Database:   driver,
		DbKey:      "neo4j",
		DBSchema:   dbschema.NewDBSchemaWithJsonSchema(nodeSchema.GetJsonSchema()),
		Env:        xtest.NewEnvConfig(),
		IsPubEvent: false,
		DaoType:    "node",
	}
	nodeDao := NewDao[map[string]any](nodeCfg)
	return nodeDao

}

func newNodes(humanName string, id string, t *testing.T) []map[string]interface{} {

	startNode := map[string]any{
		"id":         "start" + id,
		"name":       humanName,
		"analyse":    "",
		"age":        randomutils.IntMax(100),
		"peopleType": []string{"1111"},
		"tags":       []string{"tag1", "tag2"},
		"caseId":     "test",
		"graphId":    "neo4j-test",
	}
	endNode := map[string]any{
		"id":      "end" + id,
		"name":    humanName,
		"analyse": "",
		"age":     randomutils.IntMax(100),
		//"birthday":   times.NowTime(),
		"peopleType": []string{"1111"},
		"tags":       []string{"tag1", "tag2"},
		"caseId":     "test",
		"graphId":    "neo4j-test",
	}

	return []map[string]interface{}{
		startNode,
		endNode,
	}
}
