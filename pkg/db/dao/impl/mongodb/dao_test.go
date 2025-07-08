package mongodb

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/core/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/rsql"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/schema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/gp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/idutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/randomutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/xtest"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type Human struct {
	xtest.Base `bson:",inline"`
	Name       string    `gorm:"name" bson:"name"`
	Age        int64     `gorm:"age" bson:"age"`
	Analyse    string    `gorm:"analyse" bson:"analyse"`
	Birthday   time.Time `gorm:"birthday" bson:"birthday"`
	PeopleType []string  `gorm:"people_type;type:text[]" bson:"people_type"`
	Tags       []string  `gorm:"tags;type:text[]" bson:"tags"`
	Remark     string    `gorm:"remark;type:text" bson:"remark"`
}

var DB_NAME = "test"

func Test_Dao(t *testing.T) {
	xtest.InitEnv_MongoLocal()

	humanName := randomutils.NameCN()
	humanSchema := schema.NewJsonSchemaWithJson("human.json", xtest.HumanSchema)

	daoCfg := &idao.DaoConfig{
		DBKey:      "db",
		DBSchema:   dbschema.NewDBSchemaWithJsonSchema(humanSchema),
		Env:        env.GetEnv(),
		IsPubEvent: false,
	}

	dao := NewDao[map[string]any](daoCfg)
	ctx, err := restapp.NewTestContext(context.Background())
	if err != nil {
		t.Error(err)
		return
	}

	//dao.Table().Drop(ctx)
	//dao.Table().AutoMigrate(ctx)
	t.Run("dao.FindAll", func(t *testing.T) {
		res := dao.FindAll(ctx)
		if res.Error != nil {
			t.Error(res.Error)
		} else {
			t.Log("list:", res.GetData())
		}
	})

	newCount := int64(10)
	var list []map[string]any

	for i := int64(0); i < newCount; i++ {
		entity := map[string]any{
			"id":          randomutils.NewId(),
			"name":        humanName,
			"tenantId":    xtest.TenantId,
			"creatorName": xtest.CaseId,
			"analyse":     "",
			"age":         randomutils.Int64Max(100),
			"birthday":    randomutils.Date(),
			"peopleType":  []string{"1111"},
			"tags":        []string{"tag1", "tag2"},
		}
		list = append(list, entity)
	}

	t.Run("dao.CreateMany", func(t *testing.T) {
		gp.Try(func() error {
			res := dao.CreateMany(ctx, list)
			assert.Equal(t, newCount, res.RowsAffected)
			t.Log("RowsAffected:", res.RowsAffected)
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
		"age":        int64(1),
		"remark":     "remark," + id,
		"tags":       []string{"tag10", "tag20"},
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
			/*
				human["name"] = humanName + "2"
				human["birthday"] = times.NewDate()
				count := dao.Update(ctx, human)
				assert.Equal(t, int64(1), count.RowsAffected)

				human["birthday"] = times.NewTime()
				count = dao.Update(ctx, human)
				assert.Equal(t, int64(1), count.RowsAffected)
			*/

			human["name"] = humanName + "3"
			opts := idao.NewCallOptions().SetUpdateFields([]string{"name"})
			count3 := dao.Update(ctx, human, opts)
			assert.Equal(t, int64(1), count3.RowsAffected)
			return nil
		}).Catch(func(err error) {
			t.Error(err)
		})
	})

	t.Run("dao.FindById", func(t *testing.T) {
		gp.Try(func() error {
			e, err := dao.FindById(ctx, id)
			assert.Nil(t, err)
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
			paging := store.NewFindPagingQueryRequest()
			paging.PageSize = 2
			paging.IsTotalRows = true
			paging.Filter = rsql.NewBuilder().Eq("creatorName", "test").Build()
			res := dao.FindPaging(ctx, paging)
			assert.NotNil(t, res)
			assert.Equal(t, int64(11), res.GetTotalRows())
			assert.Equal(t, 2, len(res.GetData()))
			t.Log("totalRows=", res.GetTotalRows(), " count=", len(res.GetData()))
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
		findList, err := dao.FindByRSQL(ctx, fmt.Sprintf("creatorName=='%s'", "test"))
		assert.Nil(t, err)
		t.Log("list:", findList)
		assert.Equal(t, newCount, int64(len(findList)))
	})

	t.Run("dao.CountByRSQL", func(t *testing.T) {
		gp.Try(func() error {
			res, err := dao.CountByRSQL(ctx, fmt.Sprintf("creatorName=='%s'", "test"))
			assert.Nil(t, err)
			assert.Equal(t, newCount, res)
			t.Log("count:", res)
			return nil
		}).Catch(func(err error) {
			t.Error(err)
		})
	})

	t.Run("dao.Sum", func(t *testing.T) {
		gp.Try(func() error {
			var vals []*store.ValueCol
			vals = append(vals, &store.ValueCol{
				AggFunc: "sum",
				Field:   "age",
			})
			res, err := dao.SumByRSQL(ctx, "", vals)
			assert.Nil(t, err)
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

	t.Run("dao.FindByRSQL", func(t *testing.T) {})
}

func Test_DaoStruct(t *testing.T) {
	xtest.InitEnv_MongoLocal()

	humanName := randomutils.NameCN()
	//humanSchema := schema.NewJsonSchemaWithJson("human.json", xtest.HumanSchema)

	daoCfg := &idao.DaoConfig{
		DBKey:      "db",
		DBSchema:   dbschema.NewDBSchemaWithStruct("human", &Human{}, "human"),
		Env:        env.GetEnv(),
		IsPubEvent: false,
	}

	dao := NewDao[*Human](daoCfg)
	ctx, err := restapp.NewTestContext(context.Background())
	if err != nil {
		t.Error(err)
		return
	}

	newCount := int64(10)
	var list []*Human

	for i := int64(0); i < newCount; i++ {
		entity := &Human{
			Base: xtest.Base{
				Id:       randomutils.NewId(),
				TenantId: xtest.TenantId,
				CaseId:   xtest.CaseId,
			},
			Name:       humanName,
			Analyse:    "",
			Age:        randomutils.Int64Max(100),
			Birthday:   randomutils.Date(),
			PeopleType: []string{"1111"},
			Tags:       []string{"tag1", "tag2"},
		}
		list = append(list, entity)
	}

	t.Run("dao.CreateMany", func(t *testing.T) {
		gp.Try(func() error {
			res := dao.CreateMany(ctx, list)
			assert.Equal(t, newCount, res.RowsAffected)
			t.Log("RowsAffected:", res.RowsAffected)
			return nil
		}).Catch(func(err error) {
			t.Error(err)
		})
	})

	t.Run("dao.FindAll", func(t *testing.T) {
		res := dao.FindAll(ctx)
		if res.Error != nil {
			t.Error(res.Error)
		} else {
			t.Log("list:", res.GetData())
		}
	})

	t.Run("dao.UpdateMany", func(t *testing.T) {
		gp.Try(func() error {
			for _, v := range list {
				v.Remark = "remark," + randomutils.String(10)
			}
			dao.UpdateMany(ctx, list)
			return nil
		}).Catch(func(err error) {
			t.Error(err)
		})
	})

	id := idutils.NewId()
	human := &Human{
		Base: xtest.Base{
			Id:       randomutils.NewId(),
			TenantId: xtest.TenantId,
			CaseId:   xtest.CaseId,
		},
		Analyse:    "",
		Birthday:   time.Now(),
		PeopleType: []string{"1111"},
		Name:       humanName,
		Age:        int64(1),
		Remark:     "remark," + id,
		Tags:       []string{"tag10", "tag20"},
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
			/*
				human["name"] = humanName + "2"
				human["birthday"] = times.NewDate()
				count := dao.Update(ctx, human)
				assert.Equal(t, int64(1), count.RowsAffected)

				human["birthday"] = times.NewTime()
				count = dao.Update(ctx, human)
				assert.Equal(t, int64(1), count.RowsAffected)
			*/

			human.Name = humanName + "3"
			opts := idao.NewCallOptions().SetUpdateFields([]string{"name"})
			count3 := dao.Update(ctx, human, opts)
			assert.Equal(t, int64(1), count3.RowsAffected)
			return nil
		}).Catch(func(err error) {
			t.Error(err)
		})
	})

	t.Run("dao.FindById", func(t *testing.T) {
		gp.Try(func() error {
			e, err := dao.FindById(ctx, id)
			assert.NoError(t, err)
			assert.NotNil(t, e)
			assert.Equal(t, e.Id, id)
			t.Log("findById:", e)
			return nil
		}).Catch(func(err error) {
			t.Error(err)
		})
	})

	t.Run("dao.FindPaging", func(t *testing.T) {
		gp.Try(func() error {
			paging := store.NewFindPagingQueryRequest()
			paging.PageSize = 2
			paging.IsTotalRows = true
			paging.Filter = rsql.NewBuilder().Eq("creatorName", "test").Build()
			res := dao.FindPaging(ctx, paging)
			assert.NotNil(t, res)
			assert.Equal(t, int64(11), res.GetTotalRows())
			assert.Equal(t, 2, len(res.GetData()))
			t.Log(" totalRows= ", res.GetTotalRows(), " count= ", len(res.GetData()))
			return nil
		}).Catch(func(err error) {
			t.Error(err)
		})
	})

	t.Run("dao.UpdateByRSQL", func(t *testing.T) {
		gp.Try(func() error {
			humanName = humanName + "3"
			human.Name = humanName
			count := dao.UpdateByRSQL(ctx, fmt.Sprintf("id=='%s'", id), human)
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
		findList, err := dao.FindByRSQL(ctx, fmt.Sprintf("creatorName=='%s'", "test"))
		assert.Nil(t, err)
		t.Log("list:", findList)
		assert.Equal(t, newCount, int64(len(findList)))
	})

	t.Run("dao.CountByRSQL", func(t *testing.T) {
		gp.Try(func() error {
			res, err := dao.CountByRSQL(ctx, fmt.Sprintf("creatorName=='%s'", "test"))
			assert.Nil(t, err)
			assert.Equal(t, newCount, res)
			t.Log("count:", res)
			return nil
		}).Catch(func(err error) {
			t.Error(err)
		})
	})

	t.Run("dao.Sum", func(t *testing.T) {
		gp.Try(func() error {
			var vals []*store.ValueCol
			vals = append(vals, &store.ValueCol{
				AggFunc: "sum",
				Field:   "age",
			})
			res, err := dao.SumByRSQL(ctx, "", vals)
			assert.Nil(t, err)
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
				ids = append(ids, e.Id)
			}

			res2 := dao.DeleteByIds(ctx, ids)
			t.Log("DeleteByIds count:", res2.RowsAffected)
			assert.Equal(t, newCount, res2.RowsAffected)
			return nil
		}).Catch(func(err error) {
			t.Error(err)
		})
	})

	t.Run("dao.FindByRSQL", func(t *testing.T) {})
}

func newHumanDao() idao.Dao[*Human] {
	xtest.InitEnv_MongoLocal()
	daoCfg := &idao.DaoConfig{
		DBKey:      "db",
		DBSchema:   dbschema.NewDBSchemaWithStruct("human", &Human{}, "human"),
		Env:        env.GetEnv(),
		IsPubEvent: false,
	}

	return NewDao[*Human](daoCfg)
}

func Test_DaoStruct_Update(t *testing.T) {
	dao := newHumanDao()
	ctx, err := restapp.NewTestContext(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	humanName := "updateName"
	id := idutils.NewId()
	human := &Human{
		Base: xtest.Base{
			Id:       randomutils.NewId(),
			TenantId: xtest.TenantId,
			CaseId:   xtest.CaseId,
		},
		Analyse:    "",
		Birthday:   time.Now(),
		PeopleType: []string{"1111"},
		Name:       humanName,
		Age:        int64(1),
		Remark:     "remark," + id,
		Tags:       []string{"tag10", "tag20"},
	}

	res1 := dao.Create(ctx, human)
	assert.Equal(t, int64(1), res1.RowsAffected)

	human.Name = humanName + "1"
	res2 := dao.Update(ctx, human, idao.NewCallOptions().SetUpdateFields([]string{"name"}))
	assert.Equal(t, int64(1), res2.RowsAffected)

	human.Name = humanName + "2"
	res3 := dao.Update(ctx, human)
	assert.Equal(t, int64(1), res3.RowsAffected)
}
